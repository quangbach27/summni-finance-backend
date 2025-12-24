package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sumni-finance-backend/internal/common/logs"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/config"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type AuthHandlerInterface interface {
	HandleLogin(http.ResponseWriter, *http.Request)
	HandleCallback(http.ResponseWriter, *http.Request)
	HandleLogout(http.ResponseWriter, *http.Request)
	AuthMiddleware(next http.Handler) http.Handler
}

func HandleServerFromMux(r chi.Router, si AuthHandlerInterface) http.Handler {
	r.Route("/v1/auth", func(r chi.Router) {
		r.Get("/login", si.HandleLogin)
		r.Get("/callback", si.HandleCallback)
	})

	return r
}

type authHandler struct {
	config       oauth2.Config
	verifier     *oidc.IDTokenVerifier
	mu           sync.RWMutex
	sessionStore map[string]*oauth2.Token
}

func NewAuthHandler() *authHandler {
	cfg := config.GetConfig()
	keycloakCfg := cfg.Keycloak()

	maxRetry := 10
	retryInterval := 10 * time.Second

	var provider *oidc.Provider
	var err error

	for i := 0; i < maxRetry; i++ {
		// Use a background context with a timeout for each specific attempt
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		provider, err = oidc.NewProvider(ctx, keycloakCfg.RealmURL())
		cancel()

		if err == nil {
			slog.Info("Successfully connected to Keycloak provider")
			break
		}

		slog.Warn("Failed to query Keycloak provider, retrying...",
			"attempt", i+1,
			"max_retries", maxRetry,
			"error", err)

		if i < maxRetry-1 {
			time.Sleep(retryInterval)
		}
	}

	if provider == nil {
		slog.Error("Could not initialize OIDC provider after max retries. Exiting.")
		panic("auth provider initialization failed: " + err.Error())
	}

	return &authHandler{
		config: oauth2.Config{
			ClientID:     keycloakCfg.ClientID(),
			ClientSecret: keycloakCfg.ClientSecret(),
			RedirectURL:  keycloakCfg.RedirectURL(),
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
		verifier: provider.Verifier(&oidc.Config{
			ClientID: keycloakCfg.ClientID(),
		}),
	}
}

func (h *authHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	logger := logs.FromContext(r.Context())
	logger.Info("Handle Login")

	logger.Debug("set p_state cookie")
	state := uuid.New().String()
	h.SetCookie(w, "p_state", state, 300)

	// Redirect user to Keycloak
	url := h.config.AuthCodeURL(state)
	logger.Debug("Handle redirect to keycloak login page", "url", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *authHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	logger := logs.FromContext(r.Context())
	logger.Info("Handle auth callback")

	// 1. Verify State (CSRF Protection)
	logger.Info("verify p_state")
	stateCookie, err := r.Cookie("p_state")
	if err != nil || r.URL.Query().Get("state") != stateCookie.Value {
		httperr.BadRequest("invalid-state-for-auth-callback", errors.New("error state"), w, r)
		return
	}

	h.SetCookie(w, "p_state", "", -1) // Delete state cookie

	// 2. Exchange Code for Token
	logger.Info("Exchange Authorization Code for token")
	code := r.URL.Query().Get("code")

	token, err := h.config.Exchange(r.Context(), code)
	if err != nil {
		httperr.InternalError("failed-to-exchange-token", err, w, r)
	}

	// 3. Store session in-memory
	logger.Info("store sessionId in in-memory")
	sessionID := uuid.New().String()
	h.mu.Lock()
	h.sessionStore[sessionID] = token
	h.mu.Unlock()

	// 4. Set Session Cookie
	h.SetCookie(w, "p_sessionid", sessionID, 3600)

	http.Redirect(w, r, "http://localhost:3000/asset-source", http.StatusFound)
}

func (h *authHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
}

func (h *authHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logs.FromContext(r.Context())
		logger.Debug("authenticate ...")

		// 1. Check for Session Cookie
		cookie, err := r.Cookie("p_sessionid")
		if err != nil {
			httperr.Unauthorised("empty-sessionId", err, w, r)
			return
		}

		// 2. Get/Refresh Token from Store
		token, err := h.GetFreshToken(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized: Session invalid or expired", http.StatusUnauthorized)
			httperr.Unauthorised("session-invalid-or-expired", err, w, r)
			return
		}

		// 3. CRYPTOGRAPHIC VERIFICATION
		_, err = h.verifier.Verify(r.Context(), token.AccessToken)
		if err != nil {
			httperr.Unauthorised("token-verification-failed", err, w, r)
			return
		}

		// 4. Inject into Header for downstream Business Logic
		r.Header.Set("Authorization", "Bearer "+token.AccessToken)

		next.ServeHTTP(w, r)
	})
}

func (s *authHandler) GetFreshToken(ctx context.Context, sessionID string) (*oauth2.Token, error) {
	s.mu.RLock()
	token, exists := s.sessionStore[sessionID]
	s.mu.RUnlock()

	if !exists {
		return nil, errors.New("no sessionId found")
	}

	// TokenSource automatically uses refresh_token if the access_token is expired
	ts := s.config.TokenSource(ctx, token)
	freshToken, err := ts.Token()
	if err != nil {
		return nil, err
	}

	// Update store if token was refreshed
	if freshToken.AccessToken != token.AccessToken {
		s.mu.Lock()
		s.sessionStore[sessionID] = freshToken
		s.mu.Unlock()
	}

	return freshToken, nil
}

func (h *authHandler) SetCookie(w http.ResponseWriter, name string, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true for HTTPS
		SameSite: http.SameSiteLaxMode,
	})
}

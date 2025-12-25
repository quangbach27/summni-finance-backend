package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sumni-finance-backend/internal/common/logs"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/config"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

const (
	SessionCookie = "p_sessionid"
	SessionMaxAge = 3600 // 1 hour
	StateCookie   = "p_state"
	StateMaxAge   = 300 // 5 min
)

type ctxKey string

const ClaimsKey ctxKey = "user_claims"

type TokenClaims struct {
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	PreferredUser string `json:"preferred_username"`
	// Keycloak specific roles are often nested under 'realm_access'
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

type AuthHandlerInterface interface {
	HandleLogin(http.ResponseWriter, *http.Request)
	HandleCallback(http.ResponseWriter, *http.Request)
	HandleLogout(http.ResponseWriter, *http.Request)
}

func HandleServerFromMux(r chi.Router, si AuthHandlerInterface) http.Handler {
	r.Route("/v1/auth", func(r chi.Router) {
		r.Get("/login", si.HandleLogin)
		r.Get("/callback", si.HandleCallback)
		r.Get("/logout", si.HandleLogout)
	})

	return r
}

type TokenRepository interface {
	GetBySessionID(ctx context.Context, sessionID string) (*oauth2.Token, error)
	Save(ctx context.Context, sessionID string, token *oauth2.Token) error
	DeleteBySessionID(ctx context.Context, sessionID string) error
}

type authHandler struct {
	config    oauth2.Config
	verifier  *oidc.IDTokenVerifier
	tokenRepo TokenRepository
}

func NewAuthHandler(tokenRepo TokenRepository) *authHandler {
	kcConfig := config.GetConfig().Keycloak()

	maxRetry := 10
	retryInterval := 10 * time.Second

	var provider *oidc.Provider
	var err error

	for i := 0; i < maxRetry; i++ {
		// Use a background context with a timeout for each specific attempt
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		provider, err = oidc.NewProvider(ctx, kcConfig.RealmURL())
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
			ClientID:     kcConfig.ClientID(),
			ClientSecret: kcConfig.ClientSecret(),
			RedirectURL:  kcConfig.CallbackURL(),
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
		verifier: provider.Verifier(&oidc.Config{
			ClientID: kcConfig.ClientID(),
		}),
		tokenRepo: tokenRepo,
	}
}

func (handler *authHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	logger := logs.FromContext(r.Context())
	state := uuid.New().String()

	logger.Info("Login initiated", "state", state)

	handler.setCookie(w, StateCookie, state, StateMaxAge)

	// Redirect user to Keycloak
	url := handler.config.AuthCodeURL(state)
	logger.Debug("Redirecting to Keycloak authorization endpoint")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (handler *authHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	logger := logs.FromContext(r.Context())
	queryState := r.URL.Query().Get("state")

	logger.Info("OAuth callback received",
		"has_code", r.URL.Query().Get("code") != "",
		"has_state", queryState != "")

	// 1. Verify State (CSRF Protection)
	stateCookie, err := r.Cookie(StateCookie)
	if err != nil {
		logger.Warn("State cookie missing in callback", "error", err)
		httperr.BadRequest("invalid-state-for-auth-callback", errors.New("state cookie missing"), w, r)
		return
	}

	if queryState != stateCookie.Value {
		httperr.BadRequest("invalid-state-for-auth-callback", errors.New("State mismatch detected - possible CSRF attack"), w, r)
		return
	}

	handler.setCookie(w, StateCookie, "", -1) // Delete state cookie

	// 2. Exchange Code for Token
	code := r.URL.Query().Get("code")
	logger.Debug("Exchanging authorization code for tokens")

	token, err := handler.config.Exchange(r.Context(), code)
	if err != nil {
		httperr.InternalError("failed-to-exchange-token", err, w, r)
		return
	}

	// 3. Store session in repository
	sessionID := uuid.New().String()
	err = handler.tokenRepo.Save(r.Context(), sessionID, token)
	if err != nil {
		httperr.InternalError(
			"failed-to-save-session",
			fmt.Errorf("failed to save token: %w", err),
			w, r,
		)
		return
	}

	logger.Info("Authentication successful", "session_id", sessionID)

	// 4. Set Session Cookie
	postLoginURL := config.GetConfig().Keycloak().PostLoginURL()

	handler.setCookie(w, SessionCookie, sessionID, SessionMaxAge)
	http.Redirect(w, r, postLoginURL, http.StatusFound)
}

func (handler *authHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	logger := logs.FromContext(r.Context())
	logger.Info("Handle Logout")

	cookie, err := r.Cookie(SessionCookie)
	if err != nil {
		http.Redirect(w, r, "/api/v1/auth/login", http.StatusSeeOther)
		// return
	}

	// 1. Get token bundle to retrieve the ID Token for Keycloak
	token, err := handler.getFreshToken(r.Context(), cookie.Value)
	if err != nil {
		handler.setCookie(w, SessionCookie, "", -1)
		http.Redirect(w, r, "/api/v1/auth/login", http.StatusSeeOther)
		return
	}

	// 2. Extract ID Token for the 'id_token_hint'
	idTokenHint, ok := token.Extra("id_token").(string)
	if !ok {
		logger.Warn("id_token missing in session; silent logout may fail")
	}

	// 3. Delete from repository
	_ = handler.tokenRepo.DeleteBySessionID(r.Context(), cookie.Value)

	// 4. Clear local cookie
	handler.setCookie(w, SessionCookie, "", -1)

	// 5. Build Correct OIDC Logout URL
	kcConfig := config.GetConfig().Keycloak()
	logoutURL := fmt.Sprintf("%s/protocol/openid-connect/logout?id_token_hint=%s&client_id=%s",
		kcConfig.RealmURL(),
		idTokenHint,
		kcConfig.ClientID(),
	)

	http.Redirect(w, r, logoutURL, http.StatusFound)
}

func (handler *authHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logs.FromContext(r.Context())

		cookie, err := r.Cookie(SessionCookie)
		if err != nil {
			httperr.Unauthorised("missing-session", fmt.Errorf("Authentication failed: %w", err), w, r)
			return
		}

		sessionID := cookie.Value
		logger.Debug("Authenticating request")

		// 1. Get Token (and refresh if needed)
		token, err := handler.getFreshToken(r.Context(), sessionID)
		if err != nil {
			httperr.Unauthorised("invalid-or-expired-session", fmt.Errorf("Authentication failed: %w", err), w, r)
			return
		}

		// 2. Verify Signature and Expiry
		rawIDToken, exist := token.Extra("id_token").(string)
		if !exist {
			httperr.Unauthorised("missing-id-token", errors.New("id_token missing from token reponse"), w, r)
			return
		}

		idToken, err := handler.verifier.Verify(r.Context(), rawIDToken)
		if err != nil {
			httperr.Unauthorised("token-verification-failed", err, w, r)
			return
		}

		var tokenClaim TokenClaims
		if err = idToken.Claims(&tokenClaim); err != nil {
			httperr.InternalError("failed-to-parse-claims", err, w, r)
			return
		}

		logger.Debug("Authentication successful",
			"session_id", sessionID,
			"user_id", tokenClaim.Subject,
			"email", tokenClaim.Email)

		ctx := context.WithValue(r.Context(), ClaimsKey, &tokenClaim)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (handler *authHandler) getFreshToken(ctx context.Context, sessionID string) (*oauth2.Token, error) {
	token, err := handler.tokenRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get token by sessionID: %w", err)
	}

	// TokenSource automatically uses refresh_token if the access_token is expired
	ts := handler.config.TokenSource(ctx, token)
	freshToken, err := ts.Token()
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	// Update store if token was refreshed
	if freshToken.AccessToken != token.AccessToken {
		slog.Info("Token refreshed",
			"session_id", sessionID,
			"expires_in", time.Until(freshToken.Expiry).Round(time.Second))

		err = handler.tokenRepo.Save(ctx, sessionID, freshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to save refreshed token: %w", err)
		}
	}

	return freshToken, nil
}

func (handler *authHandler) setCookie(w http.ResponseWriter, name string, value string, maxAge int) {
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

func ClaimsFromContext(ctx context.Context) (*TokenClaims, error) {
	claims, ok := ctx.Value(ClaimsKey).(*TokenClaims)
	if !ok || claims == nil {
		return nil, errors.New("no claims in context")
	}

	return claims, nil
}

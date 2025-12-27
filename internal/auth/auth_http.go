package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

const (
	SessionKey    = "p_sessionid"
	SessionMaxAge = 3600 // 1 hour
	StateKey      = "p_state"
	StateMaxAge   = 300 // 5 min
)

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

type Oauth2Client interface {
	GetAuthorizationCodeURL(state string) string
	Authenticate(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error)
	GetLogoutURL(ctx context.Context, token *oauth2.Token) (string, error)
	ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)
}

type TokenRepository interface {
	GetBySessionID(ctx context.Context, sessionID string) (*oauth2.Token, error)
	Save(ctx context.Context, sessionID string, token *oauth2.Token) error
	DeleteBySessionID(ctx context.Context, sessionID string) error
}

type authHandler struct {
	tokenRepo    TokenRepository
	oauth2Client Oauth2Client
}

func NewAuthHandler(oauth2Client Oauth2Client, tokenRepo TokenRepository) *authHandler {
	return &authHandler{
		oauth2Client: oauth2Client,
		tokenRepo:    tokenRepo,
	}
}

func (handler *authHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()

	handler.setCookie(w, StateKey, state, StateMaxAge)

	url := handler.oauth2Client.GetAuthorizationCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (handler *authHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	queryState := r.URL.Query().Get("state")

	// 1. Verify State (CSRF Protection)
	stateCookie, err := r.Cookie(StateKey)
	if err != nil {
		httperr.BadRequest("missing-state-cookie", err, w, r)
		return
	}

	if queryState != stateCookie.Value {
		httperr.BadRequest("mismatch-state-detected", errors.New("state mismatch detected - possible CSRF attack"), w, r)
		return
	}

	handler.setCookie(w, StateKey, "", -1) // Delete state cookie

	// 2. Exchange Code for Token
	code := r.URL.Query().Get("code")
	token, err := handler.oauth2Client.ExchangeCode(r.Context(), code)
	if err != nil {
		httperr.InternalError("failed-to-exchange-token", err, w, r)
		return
	}

	// 3. Store session in repository
	sessionID := uuid.New().String()
	err = handler.tokenRepo.Save(r.Context(), sessionID, token)
	if err != nil {
		httperr.InternalError("failed-to-save-session", err, w, r)
		return
	}

	// 4. Set Session Cookie
	handler.setCookie(w, SessionKey, sessionID, SessionMaxAge)

	// 5. Redirect to PostLoginURL
	http.Redirect(w, r, config.GetConfig().Keycloak().PostLoginURL(), http.StatusFound)
}

func (handler *authHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(SessionKey)
	if err != nil {
		httperr.BadRequest("missing-cookie-session", err, w, r)
		return
	}

	sesisonID := sessionCookie.Value

	token, err := handler.tokenRepo.GetBySessionID(r.Context(), sesisonID)
	if err != nil {
		httperr.BadRequest("token-not-found-in-store", err, w, r)
		return
	}

	logoutURL, err := handler.oauth2Client.GetLogoutURL(r.Context(), token)
	if err != nil {
		httperr.InternalError("failed-to-get-logoutURL-from-keycloak", err, w, r)
		return
	}

	err = handler.tokenRepo.DeleteBySessionID(r.Context(), sesisonID)
	if err != nil {
		httperr.InternalError("failed-to-delete-session-in-store", err, w, r)
		return
	}

	handler.setCookie(w, SessionKey, "", -1)

	http.Redirect(w, r, logoutURL, http.StatusFound)
}

func (handler *authHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(SessionKey)
		if err != nil {
			httperr.Unauthorised("missing-session", fmt.Errorf("authentication failed: %w", err), w, r)
			return
		}
		sessionID := sessionCookie.Value

		token, err := handler.tokenRepo.GetBySessionID(r.Context(), sessionID)
		if err != nil {
			httperr.Unauthorised("fail-to-get-token-by-session-id", err, w, r)
			return
		}

		freshToken, err := handler.oauth2Client.Authenticate(r.Context(), token)
		if err != nil {
			httperr.Unauthorised("failed-to-authenticate-token", err, w, r)
			return
		}

		// If token refreshed, store it to store
		if token.AccessToken != freshToken.AccessToken {
			err = handler.tokenRepo.Save(r.Context(), sessionID, token)
			if err != nil {
				httperr.InternalError("failed-to-save-token", err, w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
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

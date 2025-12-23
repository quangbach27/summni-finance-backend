package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthServerInterface interface {
	Auth(w http.ResponseWriter, r *http.Request)
	AuthCallback(w http.ResponseWriter, r *http.Request)
}

func HandleServerFromMux(r chi.Router, si AuthServerInterface) http.Handler {
	r.Route("/v1/auth", func(r chi.Router) {
		r.Get("/", si.Auth)
		r.Get("/callback", si.AuthCallback)
	})

	return r
}

type 
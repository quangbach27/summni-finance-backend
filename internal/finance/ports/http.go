package ports

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FinanceServerInterface interface {
	CreateAssetSource(w http.ResponseWriter, r *http.Request)
	GetAssetSources(w http.ResponseWriter, r *http.Request)
}

func HandleFinanceFromMux(r chi.Router, si FinanceServerInterface) http.Handler {
	r.Route("/v1/asset-sources", func(r chi.Router) {
		r.Get("/", si.GetAssetSources)
		r.Post("/", si.CreateAssetSource)
	})

	return r
}

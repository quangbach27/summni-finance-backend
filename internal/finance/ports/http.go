package ports

import (
	"net/http"
	"sumni-finance-backend/internal/finance/app"

	"github.com/go-chi/chi/v5"
)

type FinanceServerInterface interface {
	CreateAssetSource(w http.ResponseWriter, r *http.Request)
	GetAssetSources(w http.ResponseWriter, r *http.Request)
}

type FinanceHandler struct {
	app app.Application
}

func NewFinanceServer(app app.Application) FinanceServerInterface {
	return &FinanceHandler{
		app: app,
	}
}

func HandleServerFromMux(r chi.Router, si FinanceServerInterface) http.Handler {
	r.Route("/v1/asset-sources", func(r chi.Router) {
		r.Get("/", si.GetAssetSources)
		r.Post("/", si.CreateAssetSource)
	})

	return r
}

package ports

import (
	"net/http"
	"sumni-finance-backend/internal/finance/app"

	"github.com/go-chi/chi/v5"
)

type FinanceHandlerInterface interface {
	CreateAssetSources(w http.ResponseWriter, r *http.Request)
	GetAssetSources(w http.ResponseWriter, r *http.Request)
}

type financeHandler struct {
	app app.Application
}

func NewFinanceHandler(app app.Application) *financeHandler {
	return &financeHandler{
		app: app,
	}
}

func HandleFinanceFromMux(r chi.Router, handler FinanceHandlerInterface) http.Handler {
	r.Route("/v1/asset-sources", func(r chi.Router) {
		r.Get("/", handler.GetAssetSources)
		r.Post("/", handler.CreateAssetSources)
	})

	return r
}

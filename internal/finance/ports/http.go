package ports

import (
	"net/http"
	"sumni-finance-backend/internal/finance/app"

	"github.com/go-chi/chi/v5"
)

type envelop map[string]any
type FinanceHandlerInterface interface {
	// AssetSource
	CreateAssetSources(http.ResponseWriter, *http.Request)
	GetAssetSources(http.ResponseWriter, *http.Request)

	// Wallet
	CreateWallet(http.ResponseWriter, *http.Request)
	GetAllWallets(http.ResponseWriter, *http.Request)
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

	r.Route("/v1/wallets", func(r chi.Router) {
		r.Post("/", handler.CreateWallet)
		r.Get("/", handler.GetAllWallets)
	})

	return r
}

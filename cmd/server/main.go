package main

import (
	"net/http"
	"sumni-finance-backend/internal/common/logs"
	"sumni-finance-backend/internal/common/server"
	financePorts "sumni-finance-backend/internal/finance/ports"
	financeService "sumni-finance-backend/internal/finance/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func main() {
	logs.Init()

	financeApplication := financeService.NewApplication()

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		// HealthCheck
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, map[string]string{"status": "ok"})
		})

		// Finance Port
		financePorts.HandleServerFromMux(
			router,
			financePorts.NewFinanceHandler(financeApplication),
		)

		return router
	})
}

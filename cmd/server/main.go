package main

import (
	"net/http"
	"sumni-finance-backend/internal/common/logs"
	"sumni-finance-backend/internal/common/server"
	financePorts "sumni-finance-backend/internal/finance/ports"
	"sumni-finance-backend/internal/finance/ports/handler"
	financeService "sumni-finance-backend/internal/finance/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	logs.Init()

	financeApp := financeService.NewApplication()

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		// Health
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			logs.GetLogEntry(r).Info("request id", "reqID", reqID)
			render.JSON(w, r, map[string]string{"status": "ok"})
		})

		financePorts.HandleFinanceFromMux(router, handler.NewFinanceServerInterface(financeApp))

		return router
	})
}

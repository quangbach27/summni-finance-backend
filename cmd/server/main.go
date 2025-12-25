package main

import (
	"log/slog"
	"net/http"
	"sumni-finance-backend/internal/auth"
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

		// Auth router
		// TODO: Enable authHandler when intergrate authentication
		authRepo, err := auth.NewInMemoryTokenRepository()
		if err != nil {
			slog.Error("critical failure", "err", err)
		}
		authHandler := auth.NewAuthHandler(authRepo)
		auth.HandleServerFromMux(router, authHandler)

		// Protected routes
		router.Group(func(protectedRoute chi.Router) {
			protectedRoute.Use(authHandler.AuthMiddleware)

			protectedRoute.Get("/healthp", func(w http.ResponseWriter, r *http.Request) {
				claims, _ := auth.ClaimsFromContext(r.Context())
				slog.Info("claims from token", "claims", claims)

				w.WriteHeader(http.StatusOK)
				render.JSON(w, r, map[string]string{"status": "ok"})
			})

			// Finance Port
			financePorts.HandleServerFromMux(
				protectedRoute,
				financePorts.NewFinanceHandler(financeApplication),
			)
		})

		return router
	})
}

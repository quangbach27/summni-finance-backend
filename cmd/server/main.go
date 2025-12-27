package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sumni-finance-backend/internal/auth"
	"sumni-finance-backend/internal/common/db"
	"sumni-finance-backend/internal/common/logs"
	"sumni-finance-backend/internal/common/server"
	financePorts "sumni-finance-backend/internal/finance/ports"
	financeService "sumni-finance-backend/internal/finance/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func main() {
	logs.Init()
	ctx := context.Background()

	connPool := db.MustNewPgConnectionPool(ctx)

	tokenRepo, err := auth.NewInMemoryTokenRepository()
	if err != nil {
		slog.Error("critical failure", "err", err)
		os.Exit(1)
	}

	keycloakClient, err := auth.NewKeycloakClient()
	if err != nil {
		slog.Error("critical failure", "err", err)
		os.Exit(1)
	}

	financeApplication := financeService.NewApplication(connPool)

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		// HealthCheck
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, map[string]string{"status": "ok"})
		})

		authHandler := auth.NewAuthHandler(keycloakClient, tokenRepo)
		auth.HandleServerFromMux(router, authHandler)

		// Protected routes
		router.Group(func(protectedRoute chi.Router) {
			protectedRoute.Use(authHandler.AuthMiddleware)

			protectedRoute.Get("/healthp", func(w http.ResponseWriter, r *http.Request) {
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

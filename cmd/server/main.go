package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sumni-finance-backend/internal/common/db"
	"sumni-finance-backend/internal/common/logs"
	"sumni-finance-backend/internal/common/server"
	finance_app "sumni-finance-backend/internal/finance/app"
	"sumni-finance-backend/internal/finance/ports"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func main() {
	logs.Init()
	ctx := context.Background()

	_ = db.MustNewPgConnectionPool(ctx)

	// TODO: Uncomment when enable authentication
	/*
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

		authHandler := auth.NewAuthHandler(keycloakClient, tokenRepo)
	*/

	pgPool := db.MustNewPgConnectionPool(ctx)
	financeApp, err := finance_app.NewApplication(pgPool)
	if err != nil {
		slog.Error("failed to init finance app", "error", err)
		os.Exit(1)
	}

	financeServer := ports.NewHttpServer(financeApp)

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		// HealthCheck
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, map[string]string{"status": "ok"})
		})

		// TODO: Uncomment when enable authentication
		/*
			auth.HandleServerFromMux(router, authHandler)
		*/

		// Protected routes
		router.Group(func(protectedRoute chi.Router) {
			// TODO: Uncomment when enable authentication
			/*
				protectedRoute.Use(authHandler.AuthMiddleware)
			*/
			ports.HandlerFromMux(financeServer, protectedRoute)
		})

		return router
	})
}

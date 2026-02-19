package main

import (
	"context"
	"net/http"
	"sumni-finance-backend/internal/common/db"
	"sumni-finance-backend/internal/common/logs"
	"sumni-finance-backend/internal/common/server"

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
	*/
	_ = db.MustNewPgConnectionPool(ctx)

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		// HealthCheck
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, map[string]string{"status": "ok"})
		})

		// TODO: Uncomment when enable authentication
		/*
			authHandler := auth.NewAuthHandler(keycloakClient, tokenRepo)
			auth.HandleServerFromMux(router, authHandler)
		*/

		// Protected routes
		router.Group(func(protectedRoute chi.Router) {
			// TODO: Uncomment when enable authentication
			/*
				protectedRoute.Use(authHandler.AuthMiddleware)
			*/
		})

		return router
	})
}

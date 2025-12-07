package main

import (
	"net/http"
	"sumni-finance-backend/internal/common/logs"
	"sumni-finance-backend/internal/common/server"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	logs.Init()

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			logs.GetLogEntry(r).Info("request id", "reqID", reqID)
			render.JSON(w, r, map[string]string{"status": "ok"})
		})

		return router
	})
}

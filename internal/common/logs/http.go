package logs

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type ctxLoggerKey int

const loggerKey ctxLoggerKey = 0

// Middleware factory that wraps chi's RequestLogger and adds logger to context
func Middleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Build logger with request attributes
			attrs := []any{
				"http_method", r.Method,
				"remote_addr", r.RemoteAddr,
				"uri", r.RequestURI,
			}

			if reqID := middleware.GetReqID(r.Context()); reqID != "" {
				attrs = append(attrs, "req_id", reqID)
			}

			// Create per-request logger and store it in context
			requestLogger := logger.With(attrs...)
			ctx := context.WithValue(r.Context(), loggerKey, requestLogger)
			r = r.WithContext(ctx)

			// Log request start
			requestLogger.Info("Request started")

			// Wrap writer to capture response
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()

			// Call next handler
			next.ServeHTTP(ww, r)

			// Log request completion
			elapsed := time.Since(start)
			requestLogger.Info("Request completed",
				"resp_status", ww.Status(),
				"resp_bytes_length", ww.BytesWritten(),
				"resp_elapsed", elapsed.Round(time.Millisecond/100).String(),
			)
		})
	}
}

func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}

	log, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok || log == nil {
		return slog.Default()
	}

	return log
}

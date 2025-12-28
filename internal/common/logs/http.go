package logs

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type ctxLoggerKey int

const loggerKey ctxLoggerKey = 0

// Middleware factory (same as Logrus version)
func Middleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&Logger{Logger: logger})
}

// Logger holds the base slog logger
type Logger struct {
	Logger *slog.Logger
}

func (sl *Logger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &Entry{}

	attrs := []any{
		"http_method", r.Method,
		"remote_addr", r.RemoteAddr,
		"uri", r.RequestURI,
	}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		attrs = append(attrs, "req_id", reqID)
	}

	// Create per-request logger
	entry.Logger = sl.Logger.With(attrs...)

	ctx := context.WithValue(r.Context(), loggerKey, entry.Logger)
	_ = r.WithContext(ctx)

	// Log request start
	entry.Logger.Info("Request started")

	return entry
}

// ----------------------
// Log Entry Implementation
// ----------------------

type Entry struct {
	Logger *slog.Logger
}

func (entry *Entry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	entry.Logger.Info("Request completed",
		"resp_status", status,
		"resp_bytes_length", bytes,
		"resp_elapsed", elapsed.Round(time.Millisecond/100).String(),
	)
}

func (entry *Entry) Panic(v interface{}, stack []byte) {
	entry.Logger.Error("Panic occurred",
		"panic", fmt.Sprintf("%+v", v),
		"stack", string(stack),
	)
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

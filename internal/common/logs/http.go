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
func NewStructuredLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{Logger: logger})
}

// StructuredLogger holds the base slog logger
type StructuredLogger struct {
	Logger *slog.Logger
}

func (sl *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{}

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

type StructuredLoggerEntry struct {
	Logger *slog.Logger
}

func (entry *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	entry.Logger.Info("Request completed",
		"resp_status", status,
		"resp_bytes_length", bytes,
		"resp_elapsed", elapsed.Round(time.Millisecond/100).String(),
	)
}

func (entry *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	entry.Logger.Error("Panic occurred",
		"panic", fmt.Sprintf("%+v", v),
		"stack", string(stack),
	)
}

// ----------------------
// Access Log Entry Logger (same API)
// ----------------------
func GetLogEntry(r *http.Request) *slog.Logger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
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

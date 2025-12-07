package logs

import (
	"log/slog"
	"os"
	"sumni-finance-backend/internal/config"
)

func Init() {
	var handler slog.Handler
	env := config.GetConfig().App().Env()

	if env == "dev" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

package logs

import (
	"log/slog"
	"os"
	"sumni-finance-backend/internal/config"

	"github.com/ThreeDotsLabs/humanslog"
)

func Init() {
	config := config.GetConfig()

	slogOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var logger *slog.Logger
	if config.App().Env() == "prod" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, slogOpts))
	} else {
		opts := &humanslog.Options{
			HandlerOptions:    slogOpts,
			MaxSlicePrintSize: 10,
			SortKeys:          true,
			NewLineAfterLog:   true,
			StringerFormatter: true,
			TimeFormat:        "[04:05]",
			DebugColor:        humanslog.Magenta,
		}

		logger = slog.New(humanslog.NewHandler(os.Stdout, opts))
	}

	// optional: set global logger
	slog.SetDefault(logger)
}

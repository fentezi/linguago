package logger

import (
	"log/slog"
	"os"
)

func NewLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case "dev":
		prettyHandler := NewHandler(
			&slog.HandlerOptions{
				Level:       slog.LevelDebug,
				AddSource:   true,
				ReplaceAttr: nil,
			},
		)
		logger = slog.New(prettyHandler)
	case "prod":
		logger = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				},
			),
		)
	}

	return logger
}

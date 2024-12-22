package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
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
		rotationLogFile := &lumberjack.Logger{
			Filename:   "./logs/app.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   true,
		}
		logger = slog.New(
			slog.NewJSONHandler(
				io.MultiWriter(rotationLogFile, os.Stdout),
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				},
			),
		)
	}

	return logger
}

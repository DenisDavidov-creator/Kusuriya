package logger

import (
	"log/slog"
	"os"
	"strings"
)

func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     parseLevel(os.Getenv("LOG_LEVEL")),
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo

	}
}

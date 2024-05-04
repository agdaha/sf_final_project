package logger

import (
	"log/slog"
	"os"
	"strings"
)

func New(level string, env string) *slog.Logger {

	opts := &slog.HandlerOptions{}
	l := strings.ToUpper(level)
	switch l {
	case "DEBUG":
		opts.Level = slog.LevelDebug
	case "WARN":
		opts.Level = slog.LevelWarn
	case "ERROR":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}

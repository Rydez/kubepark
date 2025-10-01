package logger

import (
	"log/slog"
	"os"
	"strings"
)

var Logger *slog.Logger

// InitLogger initializes the global logger with the specified level
func InitLogger(level string) {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	Logger = slog.New(handler)

	// Set as default logger
	slog.SetDefault(Logger)
}

// GetLogger returns the global logger instance
func GetLogger() *slog.Logger {
	if Logger == nil {
		InitLogger("info")
	}
	return Logger
}

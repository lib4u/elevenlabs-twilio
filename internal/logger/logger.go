package logger

import (
	"log/slog"
	"os"
)

func getGlobalLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
}

// Debug logs at [LevelDebug].
func Debug(msg string, args ...any) {
	getGlobalLogger().Debug(msg, args...)
}

// Warn logs at [LevelWarn].
func Warn(msg string, args ...any) {
	getGlobalLogger().Warn(msg, args...)
}

// Info logs at [LevelInfo].
func Info(msg string, args ...any) {
	getGlobalLogger().Info(msg, args...)
}

// Error logs at [LevelError].
func Error(msg string, args ...any) {
	getGlobalLogger().Error(msg, args...)
}

func Any(key string, value any) slog.Attr {
	return slog.Attr{key, slog.AnyValue(value)}
}

func String(key, value string) slog.Attr {
	return slog.Attr{key, slog.StringValue(value)}
}

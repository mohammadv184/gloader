package log

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(handler *Handler) *Logger {
	return &Logger{slog.New(handler)}
}

func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{slog.New(l.Logger.Handler().(*Handler).WithPrefix(prefix))}
}

var defaultLogger = NewLogger(NewHandler(os.Stdout, os.Stderr, "GLoader"))

func Default() *Logger {
	return defaultLogger
}

func With(args ...any) *slog.Logger {
	return defaultLogger.With(args...)
}

func WithPrefix(prefix string) *Logger {
	return defaultLogger.WithPrefix(prefix)
}

func WithGroup(name string) *slog.Logger {
	return defaultLogger.WithGroup(name)
}

func Info(msg string, args ...any) {
	Log(slog.LevelInfo, msg, args...)
}

func Debug(msg string, args ...any) {
	Log(slog.LevelDebug, msg, args...)
}

func Warn(msg string, args ...any) {
	Log(slog.LevelWarn, msg, args...)
}

func Error(msg string, args ...any) {
	Log(slog.LevelError, msg, args...)
}

func Fatal(msg string, args ...any) {
	Log(slog.LevelError, msg, args...)
	os.Exit(1)
}

func Log(level slog.Level, msg string, args ...any) {
	defaultLogger.Log(context.Background(), level, msg, args...)
}

// Attributes wrapper funcs

func String(key, value string) slog.Attr {
	return slog.String(key, value)
}

func Int64(key string, value int64) slog.Attr {
	return slog.Int64(key, value)
}

func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

func Uint64(key string, v uint64) slog.Attr {
	return slog.Uint64(key, v)
}

func Float64(key string, v float64) slog.Attr {
	return slog.Float64(key, v)
}

func Bool(key string, v bool) slog.Attr {
	return slog.Bool(key, v)
}

func Time(key string, v time.Time) slog.Attr {
	return slog.Time(key, v)
}

func Duration(key string, v time.Duration) slog.Attr {
	return slog.Duration(key, v)
}

func Group(key string, attrs ...any) slog.Attr {
	return slog.Group(key, attrs...)
}

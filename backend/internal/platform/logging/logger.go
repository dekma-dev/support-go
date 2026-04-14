package logging

import (
	"context"
	"log/slog"
	"os"

	platformhttp "support-go/backend/internal/platform/http"
)

type loggerContextKey struct{}

var defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

func SetDefault(logger *slog.Logger) {
	defaultLogger = logger
}

func Default() *slog.Logger {
	return defaultLogger
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerContextKey{}).(*slog.Logger); ok {
		return logger
	}

	if requestID := platformhttp.RequestIDFromContext(ctx); requestID != "" {
		return defaultLogger.With("request_id", requestID)
	}

	return defaultLogger
}

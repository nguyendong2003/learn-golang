package util

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

const (
	LayerApp        = "app"
	LayerBootstrap  = "bootstrap"
	LayerHTTP       = "http"
	LayerHandler    = "handler"
	LayerService    = "service"
	LayerRepository = "repository"
	LayerWorker     = "worker"
	LayerInfra      = "infra"
)

type contextLoggerKey struct{}

func SetupLoggerFromEnv() *slog.Logger {
	level := parseLogLevel(os.Getenv("LOG_LEVEL"))
	format := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT")))
	addSource := strings.EqualFold(strings.TrimSpace(os.Getenv("LOG_ADD_SOURCE")), "true")

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: addSource,
	}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	logger.Info("logger initialized",
		slog.String("layer", LayerBootstrap),
		slog.String("level", level.String()),
		slog.String("format", defaultIfEmpty(format, "text")),
		slog.Bool("add_source", addSource),
	)

	return logger
}

func Logger() *slog.Logger {
	return slog.Default()
}

func WithLayer(layer string) *slog.Logger {
	return Logger().With("layer", layer)
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		logger = Logger()
	}
	return context.WithValue(ctx, contextLoggerKey{}, logger)
}

func WithLoggerAttrs(ctx context.Context, attrs ...any) context.Context {
	logger := LoggerFromContext(ctx).With(attrs...)
	return WithLogger(ctx, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return Logger()
	}
	if logger, ok := ctx.Value(contextLoggerKey{}).(*slog.Logger); ok && logger != nil {
		return logger
	}
	return Logger()
}

func LoggerFromContextWithLayer(ctx context.Context, layer string) *slog.Logger {
	return LoggerFromContext(ctx).With("layer", layer)
}

func parseLogLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func defaultIfEmpty(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

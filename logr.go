package forge

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

type Logger = logr.Logger

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return logr.NewContext(ctx, logger)
}

func LoggerFrom(ctx context.Context) Logger {
	return logr.FromContextOrDiscard(ctx)
}

func NewLogger() Logger {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return zapr.NewLogger(zapLogger)
}

package forge

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

// Logger is an alias to logr.Logger in case
// the logging library is desired to be swapped out.
type Logger = logr.Logger

// WithLogger returns a Context from the parent Context
// with the given Logger inside of it.
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return logr.NewContext(ctx, logger)
}

// LoggerFrom returns a Logger embedded within the given Context
// or a no-op Logger if no such Logger exists.
func LoggerFrom(ctx context.Context) Logger {
	return logr.FromContextOrDiscard(ctx)
}

// NewLogger creates a new Logger.
func NewLogger() Logger {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return zapr.NewLogger(zapLogger)
}

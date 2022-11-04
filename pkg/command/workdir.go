package command

import (
	"context"
	"os"
)

type workdirKey struct{}

func WithWorkdir(ctx context.Context, workdir string) context.Context {
	return context.WithValue(ctx, workdirKey{}, workdir)
}

func WorkdirFrom(ctx context.Context) string {
	if wd, ok := ctx.Value(workdirKey{}).(string); ok {
		return wd
	}

	if wd, err := os.Getwd(); err == nil {
		return wd
	}

	return "."
}

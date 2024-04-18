package contaminate

import (
	"context"
	"io"
)

type inputKey struct{}

func WithStdin(ctx context.Context, stdin io.Reader) context.Context {
	return context.WithValue(ctx, inputKey{}, stdin)
}

func StdinFrom(ctx context.Context) io.Reader {
	stdin, _ := ctx.Value(inputKey{}).(io.Reader)
	return stdin
}

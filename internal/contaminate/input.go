package contaminate

import "context"

type inputKey struct{}

func WithInput(ctx context.Context, input []byte) context.Context {
	return context.WithValue(ctx, inputKey{}, input)
}

func InputFrom(ctx context.Context) []byte {
	input, _ := ctx.Value(inputKey{}).([]byte)
	return input
}

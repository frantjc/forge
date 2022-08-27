package contaminate

import (
	"context"

	"github.com/frantjc/forge"
)

type mountKey struct{}

func WithMounts(ctx context.Context, mounts ...*forge.Mount) context.Context {
	return context.WithValue(ctx, mountKey{}, mounts)
}

func MountsFrom(ctx context.Context) []*forge.Mount {
	mounts, _ := ctx.Value(mountKey{}).([]*forge.Mount)
	return mounts
}

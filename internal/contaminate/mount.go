package contaminate

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/go-fn"
)

type mountKey struct{}

func WithMounts(ctx context.Context, mounts ...forge.Mount) context.Context {
	return context.WithValue(ctx, mountKey{}, mounts)
}

func MountsFrom(ctx context.Context) []forge.Mount {
	mounts, _ := ctx.Value(mountKey{}).([]forge.Mount)
	return mounts
}

func OverrideWithMountsFrom(ctx context.Context, mounts ...forge.Mount) []forge.Mount {
	mountsFrom := MountsFrom(ctx)
	return append(fn.Filter(mounts, func(m forge.Mount, _ int) bool {
		return !fn.Some(mountsFrom, func(n forge.Mount, _ int) bool {
			return m.Destination == n.Destination
		})
	}), mountsFrom...)
}

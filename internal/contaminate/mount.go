package contaminate

import (
	"context"

	"github.com/frantjc/forge"
	xslice "github.com/frantjc/x/slice"
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
	return append(xslice.Filter(mounts, func(m forge.Mount, _ int) bool {
		return !xslice.Some(mountsFrom, func(n forge.Mount, _ int) bool {
			return m.Destination == n.Destination
		})
	}), mountsFrom...)
}

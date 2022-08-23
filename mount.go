package forge

import "context"

type Mount = ContainerConfig_Mount

type mountKey struct{}

func WithMounts(ctx context.Context, mounts ...*Mount) context.Context {
	return context.WithValue(ctx, mountKey{}, mounts)
}

func MountsFrom(ctx context.Context) []*Mount {
	mounts, _ := ctx.Value(mountKey{}).([]*Mount)
	return mounts
}

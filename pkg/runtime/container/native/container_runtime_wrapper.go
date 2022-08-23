package native

import (
	"context"

	"github.com/frantjc/forge"
)

func New(c forge.ContainerRuntime) *ContainerRuntime {
	return &ContainerRuntime{c}
}

type ContainerRuntime struct {
	forge.ContainerRuntime
}

func (r *ContainerRuntime) PullImage(ctx context.Context, reference string) (forge.Image, error) {
	return PullImage(ctx, reference)
}

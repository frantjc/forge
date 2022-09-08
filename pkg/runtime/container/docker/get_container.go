package docker

import (
	"context"

	"github.com/frantjc/forge"
)

func (d *ContainerRuntime) GetContainer(ctx context.Context, id string) (forge.Container, error) {
	_, err := d.ContainerInspect(ctx, id)
	return &Container{id, d.Client}, err
}

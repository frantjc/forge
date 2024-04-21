package docker

import (
	"context"

	"github.com/frantjc/forge"
)

func (d *ContainerRuntime) GetContainer(ctx context.Context, id string) (forge.Container, error) {
	if _, err := d.ContainerInspect(ctx, id); err != nil {
		return nil, err
	}

	return &Container{id, d.Client}, nil
}

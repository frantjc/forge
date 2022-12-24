package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/frantjc/forge"
)

func (c *Container) Start(ctx context.Context) error {
	_ = forge.LoggerFrom(ctx)
	return c.Client.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
}

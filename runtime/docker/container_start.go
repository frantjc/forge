package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
)

func (c *Container) Start(ctx context.Context) error {
	return c.Client.ContainerStart(ctx, c.ID, container.StartOptions{})
}

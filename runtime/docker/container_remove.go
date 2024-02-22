package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
)

func (c *Container) Remove(ctx context.Context) error {
	return c.Client.ContainerRemove(ctx, c.ID, container.RemoveOptions{
		Force: true,
	})
}

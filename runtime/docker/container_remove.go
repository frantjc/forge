package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

func (c *Container) Remove(ctx context.Context) error {
	return c.Client.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{
		Force: true,
	})
}

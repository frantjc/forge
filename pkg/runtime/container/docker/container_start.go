package docker

import (
	"context"

	"github.com/docker/docker/api/types"
)

func (c *Container) Start(ctx context.Context) error {
	return c.Client.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
}

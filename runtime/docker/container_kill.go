package docker

import (
	"context"
	"os"
)

func (c *Container) Kill(ctx context.Context) error {
	return c.Client.ContainerKill(ctx, c.ID, os.Kill.String())
}

package docker

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/container"
)

func (c *Container) Stop(ctx context.Context) error {
	seconds := -1
	if deadline, ok := ctx.Deadline(); ok {
		seconds = int(time.Until(deadline).Seconds())
	}

	return c.Client.ContainerStop(ctx, c.ID, container.StopOptions{
		Timeout: &seconds,
	})
}

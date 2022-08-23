package docker

import (
	"context"
	"time"
)

func (c *Container) Stop(ctx context.Context) error {
	timeout := time.Minute
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline)
	}

	return c.Client.ContainerStop(ctx, c.ID, &timeout)
}

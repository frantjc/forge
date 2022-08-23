package docker

import (
	"context"
	"time"
)

func (c *Container) Restart(ctx context.Context) error {
	timeout := time.Minute
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline)
	}

	return c.Client.ContainerRestart(ctx, c.ID, &timeout)
}

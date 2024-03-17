package docker

import (
	"context"
	"io"
	"path/filepath"
)

func (c *Container) CopyFrom(ctx context.Context, source string) (io.ReadCloser, error) {
	rc, _, err := c.CopyFromContainer(ctx, c.ID, filepath.Clean(source))
	return rc, err
}

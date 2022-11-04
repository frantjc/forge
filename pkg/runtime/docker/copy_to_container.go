package docker

import (
	"context"
	"io"
	"path/filepath"

	"github.com/docker/docker/api/types"
)

func (c *Container) CopyTo(ctx context.Context, destination string, content io.Reader) error {
	if rc, ok := content.(io.ReadCloser); ok {
		defer rc.Close()
	}

	return c.CopyToContainer(ctx, c.ID, filepath.Clean(destination), content, types.CopyToContainerOptions{})
}

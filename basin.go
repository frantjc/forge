package forge

import (
	"context"
	"io"
)

// Basin is a cache for Ores and their resulting Metals.
type Basin interface {
	NewReader(context.Context, string) (io.ReadCloser, error)
	NewWriter(context.Context, string) (io.WriteCloser, error)
}

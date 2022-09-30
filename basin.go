package forge

import (
	"context"
	"fmt"
	"io"
)

// Basin is a cache for Ores and their resulting Metals
type Basin interface {
	fmt.GoStringer

	NewReader(context.Context, string) (io.ReadCloser, error)
	NewWriter(context.Context, string) (io.WriteCloser, error)
}

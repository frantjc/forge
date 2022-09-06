package forge

import (
	"context"
	"fmt"
	"io"
)

// Basin is a cache for Ores and their
// resulting Metals
//
// TODO make more generic; don't interact with Ore or Metal directly
//      this is because ore.Action changes itself as it runs so it
//      does not have the same digest before and after a run
//      and it would be cleaner to get ahold of io interfaces directly
type Basin interface {
	fmt.GoStringer

	NewReader(context.Context, string) (io.ReadCloser, error)
	NewWriter(context.Context, string) (io.WriteCloser, error)
}

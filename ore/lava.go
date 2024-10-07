package ore

import (
	"context"
	"io"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
)

// Lava is an Ore representing two Ores of which the
// stdout of the first is piped to the stdin of the second.
type Lava struct {
	From forge.Ore `json:"from,omitempty"`
	To   forge.Ore `json:"to,omitempty"`
}

func (o *Lava) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (err error) {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		pw.CloseWithError(o.From.Liquify(ctx, containerRuntime, &forge.Drains{
			Out: pw,
			Err: drains.Err,
			Tty: drains.Tty,
		}))
	}()

	return o.To.Liquify(contaminate.WithStdin(ctx, pr), containerRuntime, drains)
}

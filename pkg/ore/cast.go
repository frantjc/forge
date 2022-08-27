package ore

import (
	"bytes"
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
)

type Cast struct {
	From forge.Ore `json:"from,omitempty"`
	To   *Pure     `json:"to,omitempty"`
}

func (o *Cast) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (lava *forge.Lava, err error) {
	var (
		buf = new(bytes.Buffer)
	)
	lava, err = o.From.Liquify(ctx, containerRuntime, &forge.Drains{
		Out: buf,
		Err: drains.Err,
		Tty: drains.Tty,
	})
	if err != nil {
		return
	}

	return o.To.Liquify(contaminate.WithInput(ctx, buf.Bytes()), containerRuntime, drains)
}

func (o *Cast) GetFrom() forge.Ore {
	return o.From
}

func (o *Cast) GetTo() forge.Ore {
	return o.To
}

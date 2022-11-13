package ore

import (
	"bytes"
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
)

type Lava struct {
	From forge.Ore `json:"from,omitempty"`
	To   *Pure     `json:"to,omitempty"`
}

func (o *Lava) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (metal *forge.Metal, err error) {
	buf := new(bytes.Buffer)
	metal, err = o.From.Liquify(ctx, containerRuntime, &forge.Drains{
		Out: buf,
		Err: drains.Err,
		Tty: drains.Tty,
	})
	if err != nil {
		return
	}

	return o.To.Liquify(contaminate.WithInput(ctx, buf.Bytes()), containerRuntime, drains)
}

func (o *Lava) GetFrom() forge.Ore {
	return o.From
}

func (o *Lava) GetTo() forge.Ore {
	return o.To
}

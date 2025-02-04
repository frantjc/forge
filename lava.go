package forge

import (
	"context"
	"io"
)

// Lava is an Ore representing two Ores of which the
// stdout of the first is piped to the stdin of the second.
type Lava struct {
	From Ore
	To   Ore
}

func (o *Lava) Liquify(ctx context.Context, containerRuntime ContainerRuntime, opts ...OreOpt) (err error) {
	var (
		opt    = oreOptsWithDefaults(opts...)
		pr, pw = io.Pipe()
	)

	go func() {
		defer pw.Close()

		_ = pw.CloseWithError(o.From.Liquify(ctx, containerRuntime, &OreOpts{
			Mounts: opt.Mounts,
			Streams: &Streams{
				Out:        pw,
				In:         opt.Streams.In,
				Err:        opt.Streams.Err,
				Tty:        opt.Streams.Tty,
				DetachKeys: opt.Streams.DetachKeys,
			},
		}))
	}()

	return o.To.Liquify(ctx, containerRuntime, &OreOpts{
		Mounts: opt.Mounts,
		Streams: &Streams{
			Out:        opt.Streams.Out,
			In:         pr,
			Err:        opt.Streams.Err,
			Tty:        opt.Streams.Tty,
			DetachKeys: opt.Streams.DetachKeys,
		},
	})
}

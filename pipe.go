package forge

import (
	"context"
	"io"
)

type Pipe struct {
	From Runnable
	To   Runnable
}

func (o *Pipe) Run(ctx context.Context, containerRuntime ContainerRuntime, opts ...RunOpt) (err error) {
	var (
		opt    = runOptsWithDefaults(opts...)
		pr, pw = io.Pipe()
	)

	go func() {
		defer pw.Close()

		_ = pw.CloseWithError(o.From.Run(ctx, containerRuntime, &RunOpts{
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

	return o.To.Run(ctx, containerRuntime, &RunOpts{
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

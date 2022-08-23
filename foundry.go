package forge

import (
	"context"
	"fmt"
	"io"
)

type Foundry struct {
	ContainerRuntime
	Deposit
}

func (f *Foundry) Process(ctx context.Context, o Ore, streams *Streams) (*Metal, error) {
	if f.ContainerRuntime == nil {
		return nil, fmt.Errorf("nil ContainerRuntime")
	}

	var (
		stdout = streams.Out
		stderr = streams.Err
	)
	if f.Deposit != nil {
		stdout = io.MultiWriter(stdout, io.Discard)
		stderr = io.MultiWriter(stderr, io.Discard)
	}

	lava, err := o.Liquify(ctx, f, &Streams{
		In:         streams.In,
		Out:        stdout,
		Err:        stderr,
		Tty:        streams.Tty,
		DetachKeys: streams.DetachKeys,
	})
	if err != nil {
		return nil, err
	}

	return &Metal{
		ExitCode: int(lava.GetExitCode()),
	}, nil
}

func (f *Foundry) GoString() string {
	return "&Foundry{ContainerRuntime: " + f.ContainerRuntime.GoString() + ", Deposit: " + f.Deposit.GoString() + "}"
}

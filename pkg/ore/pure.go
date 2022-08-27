package ore

import (
	"bytes"
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/contaminate"
)

func (o *Pure) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (*forge.Lava, error) {
	image, err := containerRuntime.PullImage(ctx, o.GetImage())
	if err != nil {
		return nil, err
	}

	container, err := containerRuntime.CreateContainer(ctx, image, &forge.ContainerConfig{
		Entrypoint: o.GetEntrypoint(),
		Cmd:        o.GetCmd(),
		Env:        o.GetEnv(),
		WorkingDir: forge.WorkingDir,
		Mounts:     contaminate.MountsFrom(ctx),
	})
	if err != nil {
		return nil, err
	}

	input := contaminate.InputFrom(ctx)
	if len(input) == 0 {
		input = o.Input
	}

	exitCode, err := container.Run(ctx, drains.ToStreams(bytes.NewReader(input)))
	if err != nil {
		return nil, err
	}

	return &forge.Lava{
		ExitCode: int64(exitCode),
	}, nil
}

package ore

import (
	"context"

	"github.com/frantjc/forge"
)

func (o *Pure) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, streams *forge.Streams) (*forge.Lava, error) {
	image, err := containerRuntime.PullImage(ctx, o.GetImage())
	if err != nil {
		return nil, err
	}

	container, err := containerRuntime.CreateContainer(ctx, image, &forge.ContainerConfig{
		Entrypoint: o.GetEntrypoint(),
		Cmd:        o.GetCmd(),
		Env:        o.GetEnv(),
	})
	if err != nil {
		return nil, err
	}

	exitCode, err := container.Run(ctx, streams)
	if err != nil {
		return nil, err
	}

	return &forge.Lava{
		ExitCode: int64(exitCode),
	}, nil
}

package ore

import (
	"bytes"
	"context"
	"path/filepath"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/bin"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/internal/events"
)

func (o *Pure) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, basin forge.Basin, drains *forge.Drains) (*forge.Cast, error) {
	image, err := containerRuntime.PullImage(ctx, o.GetImage())
	if err != nil {
		return nil, err
	}

	container, err := containerRuntime.CreateContainer(ctx, image, &forge.ContainerConfig{
		Entrypoint: []string{bin.ShimPath, "-s"},
		Mounts:     contaminate.MountsFrom(ctx),
	})
	if err != nil {
		return nil, err
	}
	defer container.Stop(ctx)   //nolint:errcheck
	defer container.Remove(ctx) //nolint:errcheck

	if err = container.CopyTo(ctx, filepath.Dir(bin.ShimPath), bin.NewShimTarArchive()); err != nil {
		return nil, err
	}

	if err = container.Start(ctx); err != nil {
		return nil, err
	}

	events.Emit(ctx, &events.Event{
		Type:     events.ContainerCreated.String(),
		Metadata: events.NewContainerMetadata(container),
	})

	input := contaminate.InputFrom(ctx)
	if len(input) == 0 {
		input = o.Input
	}

	exitCode, err := container.Exec(ctx, &forge.ContainerConfig{
		Entrypoint: o.GetEntrypoint(),
		Cmd:        o.GetCmd(),
		Env:        o.GetEnv(),
		WorkingDir: forge.WorkingDir,
	}, drains.ToStreams(bytes.NewReader(input)))
	if err != nil {
		return nil, err
	}

	return &forge.Cast{
		ExitCode: int64(exitCode),
	}, nil
}

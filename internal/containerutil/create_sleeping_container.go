package containerutil

import (
	"context"
	"path/filepath"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/bin"
	"github.com/frantjc/forge/internal/hooks"
)

func CreateSleepingContainer(ctx context.Context, containerRuntime forge.ContainerRuntime, image forge.Image, containerConfig *forge.ContainerConfig) (forge.Container, error) {
	_ = forge.LoggerFrom(ctx)

	container, err := containerRuntime.CreateContainer(ctx, image, &forge.ContainerConfig{
		Entrypoint: bin.ShimSleepEntrypoint,
		Mounts:     containerConfig.GetMounts(),
		Env:        containerConfig.GetEnv(),
	})
	if err != nil {
		return nil, err
	}

	if err = container.CopyTo(ctx, filepath.Dir(bin.ShimPath), bin.NewShimTarArchive()); err != nil {
		return nil, err
	}

	hooks.ContainerCreated.Dispatch(ctx, container)

	if err = container.Start(ctx); err != nil {
		return nil, err
	}

	hooks.ContainerStarted.Dispatch(ctx, container)

	return container, nil
}

package containerutil

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/bin"
	"github.com/frantjc/forge/internal/containerfs"
	"github.com/frantjc/forge/internal/hooks"
)

func CreateSleepingContainer(ctx context.Context, containerRuntime forge.ContainerRuntime, image forge.Image, containerConfig *forge.ContainerConfig) (forge.Container, error) {
	entrypoint := []string{bin.ShimPath, "sleep", "--sock", containerfs.ForgeSock}

	for _, mount := range containerConfig.Mounts {
		if mount.Source != "" && mount.Destination != "" {
			entrypoint = append(entrypoint,
				fmt.Sprintf("--mount=%s=%s", mount.Source, mount.Destination),
			)
		}
	}

	container, err := containerRuntime.CreateContainer(ctx, image, &forge.ContainerConfig{
		Entrypoint: entrypoint,
		Mounts:     containerConfig.Mounts,
		Env:        containerConfig.Env,
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

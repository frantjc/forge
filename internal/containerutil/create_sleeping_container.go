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

var NoForgeSock bool

type sockContainerWrapper struct {
	forge.Container
}

func (c *sockContainerWrapper) Exec(ctx context.Context, cc *forge.ContainerConfig, s *forge.Streams) (int, error) {
	ccc := new(forge.ContainerConfig)
	*ccc = *cc

	if NoForgeSock {
		ccc.Entrypoint = append([]string{bin.ShimPath, "exec", "--"}, ccc.Entrypoint...)
	} else {
		ccc.Entrypoint = append([]string{bin.ShimPath, "exec", fmt.Sprintf("--sock=%s", containerfs.ForgeSock), "--"}, ccc.Entrypoint...)
	}

	return c.Container.Exec(ctx, ccc, s)
}

func CreateSleepingContainer(ctx context.Context, containerRuntime forge.ContainerRuntime, image forge.Image, containerConfig *forge.ContainerConfig) (forge.Container, error) {
	entrypoint := []string{bin.ShimPath, "sleep"}

	if !NoForgeSock {
		entrypoint = append(entrypoint,
			fmt.Sprintf("--sock=%s", containerfs.ForgeSock),
		)

		for _, mount := range containerConfig.Mounts {
			if mount.Source != "" && mount.Destination != "" {
				entrypoint = append(entrypoint,
					fmt.Sprintf("--mount=%s=%s", mount.Source, mount.Destination),
				)
			}
		}
	}

	ccc := new(forge.ContainerConfig)
	*ccc = *containerConfig
	ccc.Entrypoint = entrypoint

	container, err := containerRuntime.CreateContainer(ctx, image, ccc)
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

	return &sockContainerWrapper{container}, nil
}

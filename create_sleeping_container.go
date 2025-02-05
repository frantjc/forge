package forge

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/frantjc/forge/internal/bin"
)

var (
	ShimName   = "shim"
	ScriptName = "script"
)

type SleepingShimContainer struct {
	Container
	WorkingDir          string
	InterceptDockerSock bool
}

func (c *SleepingShimContainer) Exec(ctx context.Context, cc *ContainerConfig, s *Streams) (int, error) {
	ccc := new(ContainerConfig)
	*ccc = *cc

	if c.InterceptDockerSock {
		ccc.Entrypoint = append([]string{filepath.Join(c.WorkingDir, ShimName), "exec", fmt.Sprintf("--sock=%s", InterceptingDockerSock(c.WorkingDir)), "--"}, ccc.Entrypoint...)
	} else {
		ccc.Entrypoint = append([]string{filepath.Join(c.WorkingDir, ShimName), "exec", "--"}, ccc.Entrypoint...)
	}

	return c.Container.Exec(ctx, ccc, s)
}

func createSleepingContainer(ctx context.Context, containerRuntime ContainerRuntime, image Image, containerConfig *ContainerConfig, opt *RunOpts) (Container, error) {
	entrypoint := []string{filepath.Join(opt.WorkingDir, ShimName), "sleep"}

	if opt.InterceptDockerSock {
		entrypoint = append(entrypoint,
			fmt.Sprintf("--sock=%s", InterceptingDockerSock(opt.WorkingDir)),
		)

		for _, mount := range containerConfig.Mounts {
			if mount.Source != "" && mount.Destination != "" {
				entrypoint = append(entrypoint,
					fmt.Sprintf("--mount=%s=%s", mount.Source, mount.Destination),
				)
			}
		}
	}

	ccc := new(ContainerConfig)
	*ccc = *containerConfig
	ccc.Entrypoint = entrypoint
	ccc.Cmd = nil

	container, err := containerRuntime.CreateContainer(ctx, image, ccc)
	if err != nil {
		return nil, err
	}

	if err = container.CopyTo(ctx, opt.WorkingDir, bin.NewShimTarArchive(ShimName)); err != nil {
		return nil, err
	}

	if err = container.Start(ctx); err != nil {
		return nil, err
	}

	HookContainerStarted.Dispatch(ctx, container)

	return &SleepingShimContainer{container, opt.WorkingDir, opt.InterceptDockerSock}, nil
}

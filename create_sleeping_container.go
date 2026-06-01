package forge

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/frantjc/forge/internal/bin"
	"github.com/frantjc/forge/internal/featureflags"
	"github.com/frantjc/forge/internal/hostfs"
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

	exitCode, err := c.Container.Exec(ctx, ccc, s)
	if err != nil {
		err = fmt.Errorf("sleeping container exec: %w", err)
	}
	return exitCode, err
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

	if featureflags.BindMounts {
		if err := writeShim(); err != nil {
			return nil, fmt.Errorf("write shim: %w", err)
		}
		ccc.Mounts = append(ccc.Mounts, Mount{
			Source:      shimPath,
			Destination: filepath.Join(opt.WorkingDir, ShimName),
		})
	}

	container, err := containerRuntime.CreateContainer(ctx, image, ccc)
	if err != nil {
		return nil, fmt.Errorf("create sleeping container: %w", err)
	}

	if !featureflags.BindMounts {
		if err = container.CopyTo(ctx, opt.WorkingDir, bin.NewShimTarArchive(ShimName)); err != nil {
			return nil, fmt.Errorf("copy shim to sleeping container: %w", err)
		}
	}

	if err = container.Start(ctx); err != nil {
		return nil, fmt.Errorf("start sleeping container: %w", err)
	}

	HookContainerStarted.Dispatch(ctx, container)

	return &SleepingShimContainer{container, opt.WorkingDir, opt.InterceptDockerSock}, nil
}

var (
	shimPath      = filepath.Join(hostfs.CacheHome, ShimName)
	writeShimOnce sync.Once
)

func writeShim() (err error) {
	writeShimOnce.Do(func() {
		if err = os.MkdirAll(filepath.Dir(shimPath), 0o755); err != nil {
			return
		}

		var f *os.File
		f, err = os.OpenFile(shimPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			return
		}
		defer f.Close()

		_, err = io.Copy(f, bin.NewShimReader())
	})
	return
}

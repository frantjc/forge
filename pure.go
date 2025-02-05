package forge

import (
	"context"

	xos "github.com/frantjc/x/os"
)

type Pure struct {
	Image      string
	Entrypoint []string
	Cmd        []string
	Env        []string
}

func (o *Pure) Run(ctx context.Context, containerRuntime ContainerRuntime, opts ...RunOpt) error {
	opt := runOptsWithDefaults(opts...)

	image, err := containerRuntime.PullImage(ctx, o.Image)
	if err != nil {
		return err
	}

	containerConfig := &ContainerConfig{
		Entrypoint: o.Entrypoint,
		Cmd:        o.Cmd,
		Env:        o.Env,
		WorkingDir: opt.WorkingDir,
		Mounts:     opt.Mounts,
	}

	container, err := createSleepingContainer(ctx, containerRuntime, image, containerConfig, opt)
	if err != nil {
		return err
	}
	defer container.Stop(ctx)   //nolint:errcheck
	defer container.Remove(ctx) //nolint:errcheck

	if exitCode, err := container.Exec(ctx, containerConfig, opt.Streams); err != nil {
		return err
	} else if exitCode > 0 {
		return xos.NewExitCodeError(ErrContainerExitedWithNonzeroExitCode, exitCode)
	}

	return nil
}

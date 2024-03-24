package ore

import (
	"context"
	"strings"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/cloudbuild"
	"github.com/frantjc/forge/forgecloudbuild"
	"github.com/frantjc/forge/internal/containerutil"
	"github.com/frantjc/forge/internal/contaminate"
	xos "github.com/frantjc/x/os"
)

type CloudBuild struct {
	cloudbuild.Step `json:",inline"`
}

func (o *CloudBuild) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) error {
	image, err := forgecloudbuild.GetImageForStep(ctx, containerRuntime, &o.Step)
	if err != nil {
		return err
	}

	home := "/root"
	if config, err := image.Config(); err == nil {
		for _, env := range config.Env {
			if strings.HasPrefix(env, "HOME=") {
				home = strings.TrimPrefix(env, "HOME=")
				break
			}
		}
	}

	containerConfig, script, err := forgecloudbuild.StepToContainerConfigAndScript(&o.Step, home)
	if err != nil {
		return err
	}
	containerConfig.Mounts = contaminate.OverrideWithMountsFrom(ctx, containerConfig.Mounts...)

	container, err := containerutil.CreateSleepingContainer(ctx, containerRuntime, image, containerConfig)
	if err != nil {
		return err
	}
	defer container.Stop(ctx)   //nolint:errcheck
	defer container.Remove(ctx) //nolint:errcheck

	if err = forgecloudbuild.CopyScriptToContainer(ctx, container, script); err != nil {
		return err
	}

	if exitCode, err := container.Exec(ctx, containerConfig, drains.ToStreams(nil)); err != nil {
		return err
	} else if exitCode > 0 {
		return xos.NewExitCodeError(ErrContainerExitedWithNonzeroExitCode, exitCode)
	}

	return nil
}

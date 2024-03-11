package ore

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/azuredevops"
	"github.com/frantjc/forge/forgeazure"
	"github.com/frantjc/forge/internal/containerutil"
	"github.com/frantjc/forge/internal/contaminate"
	xos "github.com/frantjc/x/os"
)

type Task struct {
	Task      string            `json:"task,omitempty"`
	Inputs    map[string]string `json:"inputs,omitempty"`
	Execution string            `json:"execution,omitempty"`
}

func (o *Task) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) error {
	ref, err := azuredevops.Parse(o.Task)
	if err != nil {
		return err
	}

	task, err := azuredevops.GetReferenceTask(ref)
	if err != nil {
		return err
	}

	execution, err := forgeazure.ParseExecution(o.Execution)
	if err != nil {
		return err
	}

	containerConfig, err := forgeazure.TaskToContainerConfig(ref, task, execution, o.Inputs)
	if err != nil {
		return err
	}
	containerConfig.Mounts = contaminate.OverrideWithMountsFrom(ctx, containerConfig.Mounts...)

	image, err := forgeazure.GetImageForExecution(ctx, containerRuntime, execution)
	if err != nil {
		return err
	}

	container, err := containerutil.CreateSleepingContainer(ctx, containerRuntime, image, containerConfig)
	if err != nil {
		return err
	}
	defer container.Stop(ctx)   //nolint:errcheck
	defer container.Remove(ctx) //nolint:errcheck

	if exitCode, err := container.Exec(ctx, containerConfig, drains.ToStreams(nil)); err != nil {
		return err
	} else if exitCode > 0 {
		return xos.NewExitCodeError(ErrContainerExitedWithNonzeroExitCode, exitCode)
	}

	return nil
}

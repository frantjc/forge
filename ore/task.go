package ore

import (
	"context"
	"path/filepath"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/azuredevops"
	"github.com/frantjc/forge/forgeazure"
	"github.com/frantjc/forge/internal/containerutil"
	errorcode "github.com/frantjc/go-error-code"
)

type Task struct {
	Task   string            `json:"task,omitempty"`
	Inputs map[string]string `json:"inputs,omitempty"`
}

func (o *Task) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) error {
	log := forge.LoggerFrom(ctx)

	ref, err := azuredevops.Parse(o.Task)
	if err != nil {
		return err
	}

	task, err := azuredevops.GetReferenceTask(ref)
	if err != nil {
		return err
	}

	image, err := forgeazure.GetImageForTask(ctx, containerRuntime, task)
	if err != nil {
		return err
	}

	taskDir, err := forgeazure.TaskReferenceToDirectory(ref)
	if err != nil {
		return err
	}

	containerConfig := &forge.ContainerConfig{
		Entrypoint: []string{"node", filepath.Join(forgeazure.DefaultMapping.TaskPath, task.Executions.Node.Target)},
		Mounts: []forge.Mount{
			{
				Source:      taskDir,
				Destination: forgeazure.DefaultMapping.TaskPath,
			},
		},
	}

	container, err := containerutil.CreateSleepingContainer(ctx, containerRuntime, image, containerConfig)
	if err != nil {
		return err
	}

	log.Info("exec", "src", taskDir, "dst", forgeazure.DefaultMapping.TaskPath)

	if exitCode, err := container.Exec(ctx, containerConfig, drains.ToStreams(nil)); err != nil {
		return err
	} else if exitCode > 0 {
		return errorcode.New(ErrContainerExitedWithNonzeroExitCode, errorcode.WithExitCode(exitCode))
	}

	if err = container.Stop(ctx); err != nil {
		return err
	}

	return container.Remove(ctx)
}

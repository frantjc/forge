package forge

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/frantjc/forge/azuredevops"
	xos "github.com/frantjc/x/os"
)

type Task struct {
	Task      string
	Inputs    map[string]string
	Execution string
}

func (o *Task) Run(ctx context.Context, containerRuntime ContainerRuntime, opts ...RunOpt) error {
	opt := runOptsWithDefaults(opts...)

	ref, err := azuredevops.Parse(o.Task)
	if err != nil {
		return err
	}

	task, err := azuredevops.GetReferenceTask(ref)
	if err != nil {
		return err
	}

	containerConfig, err := taskToContainerConfig(ref, task, o.Execution, o.Inputs, opt)
	if err != nil {
		return err
	}
	containerConfig.Mounts = overrideMounts(containerConfig.Mounts, opt.Mounts...)

	image, err := pullImageForExecution(ctx, containerRuntime, o.Execution)
	if err != nil {
		return err
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

func taskReferenceToDirectory(ref *azuredevops.TaskReference) (string, error) {
	if ref.IsLocal() {
		return filepath.Abs(ref.Path)
	}

	return "", fmt.Errorf("remote Azure DevOps tasks are not implemented")
}

func taskToContainerConfig(ref *azuredevops.TaskReference, task *azuredevops.Task, execution string, inputs map[string]string, opt *RunOpts) (*ContainerConfig, error) {
	taskDir, err := taskReferenceToDirectory(ref)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, fmt.Errorf("nil task")
	}

	var (
		target           string
		workingDirectory string
		env              = make([]string, len(task.Inputs))
	)
	if task.Executions != nil {
		switch execution {
		case azuredevops.ExecutionNode:
			target = task.Executions.Node.Target
			workingDirectory = task.Executions.Node.WorkingDirectory
		case azuredevops.ExecutionNode16:
			target = task.Executions.Node16.Target
			workingDirectory = task.Executions.Node16.WorkingDirectory
		case azuredevops.ExecutionNode10:
			target = task.Executions.Node10.Target
			workingDirectory = task.Executions.Node10.WorkingDirectory
		default:
			return nil, fmt.Errorf("powershell task not supported")
		}
	} else {
		return nil, fmt.Errorf("task has no executions")
	}

	if target == "" {
		return nil, fmt.Errorf("exeuction has no target")
	}

	for i, input := range task.Inputs {
		name := "INPUT_" + regexp.MustCompile("[ .]").ReplaceAllString(task.Name, "_")
		if in, ok := inputs[input.Name]; ok {
			env[i] = fmt.Sprintf("%s=%s", name, in)
		} else if input.Required {
			return nil, fmt.Errorf("required input %s not supplied: %s", input.Name, input.HelpMarkDown)
		} else {
			env[i] = fmt.Sprintf("%s=%s", name, input.DefaultValue)
		}
	}

	return &ContainerConfig{
		Entrypoint: []string{"node", filepath.Join(AzureDevOpsTaskWorkingDir(opt.WorkingDir), target)},
		Mounts: []Mount{
			{
				Source:      taskDir,
				Destination: AzureDevOpsTaskWorkingDir(opt.WorkingDir),
			},
		},
		Env:        env,
		WorkingDir: filepath.Join(AzureDevOpsTaskWorkingDir(opt.WorkingDir), workingDirectory),
	}, nil
}

func pullImageForExecution(ctx context.Context, containerRuntime ContainerRuntime, execution string) (Image, error) {
	ref := ""

	switch execution {
	case azuredevops.ExecutionNode:
		ref = NodeImageReference
	case azuredevops.ExecutionNode16:
		ref = Node16ImageReference
	case azuredevops.ExecutionNode10:
		ref = Node10ImageReference
	default:
		return nil, fmt.Errorf("powershell tasks unsupported")
	}

	return containerRuntime.PullImage(ctx, ref)
}

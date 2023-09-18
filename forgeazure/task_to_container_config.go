package forgeazure

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/azuredevops"
)

func TaskToContainerConfig(ref *azuredevops.TaskReference, task *azuredevops.Task, execution Execution, inputs map[string]string) (*forge.ContainerConfig, error) {
	return DefaultMapping.TaskToContainerConfig(ref, task, execution, inputs)
}

var r = regexp.MustCompile("[ .]")

func (m *Mapping) TaskToContainerConfig(ref *azuredevops.TaskReference, task *azuredevops.Task, execution Execution, inputs map[string]string) (*forge.ContainerConfig, error) {
	taskDir, err := m.TaskReferenceToDirectory(ref)
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
		case ExecutionNode:
			target = task.Executions.Node.Target
			workingDirectory = task.Executions.Node.WorkingDirectory
		case ExecutionNode16:
			target = task.Executions.Node16.Target
			workingDirectory = task.Executions.Node16.WorkingDirectory
		case ExecutionNode10:
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
		name := "INPUT_" + r.ReplaceAllString(task.Name, "_")
		if in, ok := inputs[input.Name]; ok {
			env[i] = fmt.Sprintf("%s=%s", name, in)
		} else if input.Required {
			return nil, fmt.Errorf("required input %s not supplied: %s", input.Name, input.HelpMarkDown)
		} else {
			env[i] = fmt.Sprintf("%s=%s", name, input.DefaultValue)
		}
	}

	return &forge.ContainerConfig{
		Entrypoint: []string{"node", filepath.Join(m.TaskPath, target)},
		Mounts: []forge.Mount{
			{
				Source:      taskDir,
				Destination: m.TaskPath,
			},
		},
		Env:        env,
		WorkingDir: filepath.Join(m.TaskPath, workingDirectory),
	}, nil
}

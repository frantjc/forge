package forgeazure

import (
	"context"
	"fmt"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/azuredevops"
)

const (
	// DefaultNode10ImageReference is the default image to use
	// when a task specifies a "Node10" execution.
	DefaultNode10ImageReference = "docker.io/library/node:10"
	// DefaultNode16ImageReference is the default image to use
	// when a task specifies a "Node16" execution.
	DefaultNode16ImageReference = "docker.io/library/node:16"
	// DefaultNodeImageReference is the default image to use
	// when a task specifies a "Node" execution.
	DefaultNodeImageReference = DefaultNode16ImageReference
)

var (
	// Node10ImageReference is the image to use
	// when a task specifies a "Node10" execution.
	Node10ImageReference = DefaultNode10ImageReference
	// Node16ImageReference is the image to use
	// when a task specifies a "Node16" execution.
	Node16ImageReference = DefaultNode16ImageReference
	// NodeImageReference is the image to use
	// when a task specifies a "Node" execution.
	NodeImageReference = DefaultNodeImageReference
)

// GetImageForTask ...
func GetImageForTask(ctx context.Context, containerRuntime forge.ContainerRuntime, task *azuredevops.Task) (forge.Image, error) {
	if task != nil && task.Executions != nil {
		if task.Executions.Node != nil {
			return containerRuntime.PullImage(ctx, NodeImageReference)
		} else if task.Executions.Node16 != nil {
			return containerRuntime.PullImage(ctx, Node16ImageReference)
		} else if task.Executions.Node10 != nil {
			return containerRuntime.PullImage(ctx, Node10ImageReference)
		} else if task.Executions.PowerShell != nil || task.Executions.PowerShell3 != nil {
			return nil, fmt.Errorf("powershell tasks unsupported")
		}
	}

	return nil, fmt.Errorf("task specified no execution")
}

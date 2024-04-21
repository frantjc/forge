package forgeazure

import (
	"context"
	"fmt"

	"github.com/frantjc/forge"
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

// GetImageForExecution ...
func GetImageForExecution(ctx context.Context, containerRuntime forge.ContainerRuntime, execution Execution) (forge.Image, error) {
	ref := ""

	switch execution {
	case ExecutionNode:
		ref = NodeImageReference
	case ExecutionNode16:
		ref = Node16ImageReference
	case ExecutionNode10:
		ref = Node10ImageReference
	default:
		return nil, fmt.Errorf("powershell tasks unsupported")
	}

	return containerRuntime.PullImage(ctx, ref)
}

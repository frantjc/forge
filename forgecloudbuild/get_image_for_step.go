package forgecloudbuild

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/cloudbuild"
)

func GetImageForStep(ctx context.Context, containerRuntime forge.ContainerRuntime, step *cloudbuild.Step) (forge.Image, error) {
	return containerRuntime.PullImage(ctx, step.Name)
}

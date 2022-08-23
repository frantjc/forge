package actions2container

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/github/actions"
)

func CreateVolumes(ctx context.Context, containerRuntime forge.ContainerRuntime, uses *actions.Uses) ([]forge.Volume, error) {
	actionVolume, err := containerRuntime.CreateVolume(ctx, UsesToVolumeName(uses))
	if err != nil {
		return nil, err
	}

	runnerToolCacheVolume, err := containerRuntime.CreateVolume(ctx, DefaultRunnerToolCacheVolumeName)
	if err != nil {
		return nil, err
	}

	return []forge.Volume{actionVolume, runnerToolCacheVolume}, nil
}

package actions2container

import (
	"context"
	"path/filepath"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/bin"
	"github.com/frantjc/forge/internal/contaminate"
	"github.com/frantjc/forge/pkg/github/actions"
)

func CreateContainer(ctx context.Context, containerRuntime forge.ContainerRuntime, image forge.Image, containerConfig *forge.ContainerConfig) (forge.Container, error) {
	container, err := containerRuntime.CreateContainer(ctx, image, containerConfig)
	if err != nil {
		return nil, err
	}

	return container, container.CopyTo(ctx, filepath.Dir(bin.ShimPath), bin.NewShimTarArchive())
}

func CreateContainerForUses(ctx context.Context, containerRuntime forge.ContainerRuntime, uses *actions.Uses) (forge.Container, error) {
	containerConfig := UsesToConfig(uses)
	containerConfig.Mounts = append(containerConfig.Mounts, contaminate.MountsFrom(ctx)...)

	return CreateContainer(ctx, containerRuntime, UsesImage, UsesToConfig(uses))
}

package actions2container

import (
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/bin"
	"github.com/frantjc/forge/pkg/github/actions"
)

var (
	UsesImage = bin.NewShimImage()
)

func UsesToConfig(uses *actions.Uses) *forge.ContainerConfig {
	return DefaultMap.UsesToConfig(uses)
}

func (m *Map) UsesToConfig(uses *actions.Uses) *forge.ContainerConfig {
	return &forge.ContainerConfig{
		Entrypoint: append(bin.ShimEntrypoint, "-c", uses.String(), m.ActionPath),
		WorkingDir: forge.WorkingDir,
		Mounts: []*forge.Mount{
			{
				Source:      UsesToVolumeName(uses),
				Destination: m.ActionPath,
			},
		},
	}
}

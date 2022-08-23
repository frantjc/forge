package actions2container

import (
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/github/actions"
)

func UsesToMount(uses *actions.Uses) *forge.ContainerConfig_Mount {
	return DefaultMap.UsesToMount(uses)
}

func (m *Map) UsesToMount(uses *actions.Uses) *forge.ContainerConfig_Mount {
	return &forge.ContainerConfig_Mount{
		Source:      UsesToVolumeName(uses),
		Destination: m.ActionPath,
	}
}

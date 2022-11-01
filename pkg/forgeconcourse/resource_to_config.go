package forgeconcourse

import (
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/concourse"
)

func ResourceToConfig(resource *concourse.Resource, resourceType *concourse.ResourceType, method string) *forge.ContainerConfig {
	return DefaultMapping.ResourceToConfig(resource, resourceType, method)
}

func (m *Mapping) ResourceToConfig(resource *concourse.Resource, resourceType *concourse.ResourceType, method string) *forge.ContainerConfig {
	return &forge.ContainerConfig{
		Entrypoint: GetEntrypoint(method),
		Cmd:        []string{m.GetRootPath() + "/" + resource.GetName()},
		Privileged: resourceType.GetPrivileged(),
		Mounts: []*forge.Mount{
			{
				Destination: m.GetRootPath() + "/" + resource.GetName(),
			},
		},
	}
}

const (
	MethodGet = concourse.MethodGet
	MethodPut = concourse.MethodPut
)

const (
	EntrypointGet = "/opt/resource/in"
	EntrypointPut = "/opt/resource/out"
)

func GetEntrypoint(method string) []string {
	switch method {
	case MethodGet:
		return []string{EntrypointGet}
	case MethodPut:
		return []string{EntrypointPut}
	}

	return nil
}

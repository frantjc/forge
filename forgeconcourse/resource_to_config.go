package forgeconcourse

import (
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/concourse"
)

func ResourceToConfig(resource *concourse.Resource, resourceType *concourse.ResourceType, method string) *forge.ContainerConfig {
	return DefaultMapping.ResourceToConfig(resource, resourceType, method)
}

func (m *Mapping) ResourceToConfig(resource *concourse.Resource, resourceType *concourse.ResourceType, method string) *forge.ContainerConfig {
	return &forge.ContainerConfig{
		Entrypoint: GetEntrypoint(method),
		Cmd:        []string{m.RootPath + "/" + resource.Name},
		Privileged: resourceType.Privileged,
		Mounts: []*forge.Mount{
			{
				Destination: m.RootPath + "/" + resource.Name,
			},
		},
	}
}

const (
	MethodGet   = concourse.MethodGet
	MethodPut   = concourse.MethodPut
	MethodCheck = concourse.MethodCheck
)

const (
	EntrypointGet   = "/opt/resource/in"
	EntrypointPut   = "/opt/resource/out"
	EntrypointCheck = "/opt/resource/check"
)

func GetEntrypoint(method string) []string {
	switch method {
	case MethodGet:
		return []string{EntrypointGet}
	case MethodPut:
		return []string{EntrypointPut}
	case MethodCheck:
		return []string{EntrypointCheck}
	}

	return nil
}

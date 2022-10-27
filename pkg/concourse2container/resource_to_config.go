package concourse2container

import (
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/concourse"
)

const (
	DefaultRootPath = forge.WorkingDir
)

func ResourceToConfig(resource *concourse.Resource, resourceType *concourse.ResourceType, method string) *forge.ContainerConfig {
	return &forge.ContainerConfig{
		Entrypoint: ResourceMethod(method).Entrypoint(),
		Cmd:        []string{DefaultRootPath + "/" + resource.GetName()},
		Privileged: resourceType.GetPrivileged(),
		Mounts: []*forge.Mount{
			{
				Destination: DefaultRootPath + "/" + resource.GetName(),
			},
		},
	}
}

type ResourceMethod string

const (
	ResourceMethodGet ResourceMethod = "get"
	ResourceMethodPut ResourceMethod = "put"
)

const (
	EntrypointGet = "/opt/resource/in"
	EntrypointPut = "/opt/resource/out"
)

func (m ResourceMethod) Entrypoint() []string {
	switch m {
	case ResourceMethodGet:
		return []string{EntrypointGet}
	case ResourceMethodPut:
		return []string{EntrypointPut}
	}

	return nil
}

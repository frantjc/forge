package docker

import "github.com/docker/docker/client"

func New(c *client.Client) *ContainerRuntime {
	return &ContainerRuntime{c}
}

type ContainerRuntime struct {
	*client.Client
}

func (f *ContainerRuntime) GoString() string {
	return "&ContainerRuntime{" + f.DaemonHost() + "}"
}

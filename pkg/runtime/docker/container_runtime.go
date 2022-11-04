package docker

import (
	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
)

func New(c *client.Client) forge.ContainerRuntime {
	return &ContainerRuntime{c}
}

type ContainerRuntime struct {
	*client.Client
}

func (f *ContainerRuntime) GoString() string {
	return "&ContainerRuntime{" + f.DaemonHost() + "}"
}

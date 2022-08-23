package docker

import (
	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/runtime/container/native"
)

func New(c *client.Client) forge.ContainerRuntime {
	return &ContainerRuntime{c, &native.ContainerRuntime{}}
}

type ContainerRuntime struct {
	*client.Client
	*native.ContainerRuntime
}

func (f *ContainerRuntime) GoString() string {
	return "&ContainerRuntime{" + f.DaemonHost() + "}"
}

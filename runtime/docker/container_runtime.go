package docker

import "github.com/docker/docker/client"

func New(c *client.Client, dind bool) *ContainerRuntime {
	return &ContainerRuntime{c, dind}
}

type ContainerRuntime struct {
	*client.Client
	DockerInDocker bool
}

func (f *ContainerRuntime) GoString() string {
	return "&ContainerRuntime{" + f.DaemonHost() + "}"
}

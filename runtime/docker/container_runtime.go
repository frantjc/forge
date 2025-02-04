package docker

import "github.com/docker/docker/client"

func New(c *client.Client, dindPath string) *ContainerRuntime {
	return &ContainerRuntime{c, dindPath}
}

// ContainerRuntime implements github.com/frantjc/forge.ContainerRuntime.
type ContainerRuntime struct {
	// Client interacts with a Docker daemon.
	*client.Client
	// DockerInDockerPath signals whether or not to mount the docker.sock of the
	// *github.com/docker/docker/client.Client and configuration to direct
	// `docker` to it into each container that it runs.
	DockerInDockerPath string
}

func (f *ContainerRuntime) GoString() string {
	return "&ContainerRuntime{" + f.DaemonHost() + "}"
}

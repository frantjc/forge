package runtime

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/runtime/docker"
	"github.com/frantjc/forge/runtime/dockerd"
)

// New returns a ContainerRuntime by attempting to connect to a Docker daemon
// via the environment. If that fails, it falls back to looking for a docker,
// podman, or nerdctl binary on the PATH.
func New(ctx context.Context, dindPath string) (forge.ContainerRuntime, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err == nil {
		if _, err = cli.Ping(ctx); err == nil {
			return dockerd.New(cli, dindPath), nil
		}
	}
	return docker.New("")
}

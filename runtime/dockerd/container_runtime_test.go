//go:build dockerd

package dockerd_test

import (
	"testing"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge/runtime"
	"github.com/frantjc/forge/runtime/dockerd"
	"github.com/stretchr/testify/require"
)

func TestContainerRuntime(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	require.NoError(t, err)
	runtime.TestContainerRuntimeConformance(t, dockerd.New(cli, ""))
}

//go:build docker

package docker_test

import (
	"testing"

	"github.com/frantjc/forge/runtime"
	"github.com/frantjc/forge/runtime/docker"
	"github.com/stretchr/testify/require"
)

func TestContainerRuntime(t *testing.T) {
	cr, err := docker.New("")
	require.NoError(t, err)
	runtime.TestContainerRuntimeConformance(t, cr)
}

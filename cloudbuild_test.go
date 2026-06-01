//go:build shim && (docker || dockerd)

package forge_test

import (
	"testing"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/cloudbuild"
	"github.com/stretchr/testify/require"
)

func TestCloudBuildRun(t *testing.T) {
	cr := Runtime(t)

	step := &forge.CloudBuild{
		Step: cloudbuild.Step{
			Name:       "alpine",
			Entrypoint: "sh",
			Args:       []string{"-c", "exit 0"},
		},
	}

	require.NoError(t, step.Run(t.Context(), cr, forge.WithMountShim(), forge.WithStreams(Streams(t))))
}

func TestCloudBuildRunScript(t *testing.T) {
	ctx := t.Context()
	cr := Runtime(t)

	step := &forge.CloudBuild{
		Step: cloudbuild.Step{
			Name:   "alpine",
			Script: "#!/bin/sh\nexit 0\n",
		},
	}

	require.NoError(t, step.Run(ctx, cr, forge.WithMountShim(), forge.WithStreams(Streams(t))))
}

//go:build shim && (docker || dockerd)

package forge_test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/cloudbuild"
	xos "github.com/frantjc/x/os"
	"github.com/stretchr/testify/require"
)

func TestCloudBuildRun(t *testing.T) {
	cr := Runtime(t)

	step := &forge.CloudBuild{
		Step: cloudbuild.Step{
			Name:       "public.ecr.aws/docker/library/alpine",
			Entrypoint: "sh",
			Args:       []string{"-c", "exit 0"},
		},
	}

	require.NoError(t, step.Run(t.Context(), cr, forge.WithStreams(Streams(t)), forge.WithMountShim()))
}

func TestCloudBuildRunScript(t *testing.T) {
	ctx := t.Context()
	cr := Runtime(t)

	step := &forge.CloudBuild{
		Step: cloudbuild.Step{
			Name:   "public.ecr.aws/docker/library/alpine",
			Script: "#!/bin/sh\nexit 0\n",
		},
	}

	require.NoError(t, step.Run(ctx, cr, forge.WithStreams(Streams(t)), forge.WithMountShim()))
}

func TestCloudBuildRunNonzeroExitCode(t *testing.T) {
	cr := Runtime(t)

	expected := rand.IntN(254) + 1
	step := &forge.CloudBuild{
		Step: cloudbuild.Step{
			Name:       "public.ecr.aws/docker/library/alpine",
			Entrypoint: "sh",
			Args:       []string{"-c", fmt.Sprintf("exit %d", expected)},
		},
	}

	err := step.Run(t.Context(), cr, forge.WithStreams(Streams(t)), forge.WithMountShim())
	require.Error(t, err)
	actual := xos.ErrorExitCode(err)
	require.Equal(t, expected, actual)
}

func TestCloudBuildRunEnv(t *testing.T) {
	cr := Runtime(t)

	step := &forge.CloudBuild{
		Step: cloudbuild.Step{
			Name:       "public.ecr.aws/docker/library/alpine",
			Entrypoint: "sh",
			Args:       []string{"-c", `[ "$MY_VAR" = "hello" ]`},
			Env:        []string{"MY_VAR=hello"},
		},
	}

	require.NoError(t, step.Run(t.Context(), cr, forge.WithStreams(Streams(t)), forge.WithMountShim()))
}

func TestCloudBuildRunSubstitutions(t *testing.T) {
	cr := Runtime(t)

	step := &forge.CloudBuild{
		Step: cloudbuild.Step{
			Name:       "public.ecr.aws/docker/library/alpine",
			Entrypoint: "sh",
			Args:       []string{"-c", `[ "$_MY_SUB" = "world" ]`},
			Substitutions: map[string]string{
				"_MY_SUB": "world",
			},
			AutomapSubstitutions: true,
		},
	}

	require.NoError(t, step.Run(t.Context(), cr, forge.WithStreams(Streams(t)), forge.WithMountShim()))
}

func TestCloudBuildRunDynamicSubstitutions(t *testing.T) {
	cr := Runtime(t)

	step := &forge.CloudBuild{
		Step: cloudbuild.Step{
			Name:       "public.ecr.aws/docker/library/alpine",
			Entrypoint: "sh",
			Args:       []string{"-c", `[ "$_GREETING" = "hello world" ]`},
			Substitutions: map[string]string{
				"_WORD":     "world",
				"_GREETING": "hello $_WORD",
			},
			DynamicSubstitutions: true,
			AutomapSubstitutions: true,
		},
	}

	require.NoError(t, step.Run(t.Context(), cr, forge.WithStreams(Streams(t)), forge.WithMountShim()))
}

func TestCloudBuildRunScriptWithArgs(t *testing.T) {
	cr := Runtime(t)

	step := &forge.CloudBuild{
		Step: cloudbuild.Step{
			Name:   "public.ecr.aws/docker/library/alpine",
			Script: "#!/bin/sh\n[ \"$1\" = \"hello\" ]\n",
			Args:   []string{"hello"},
		},
	}

	require.NoError(t, step.Run(t.Context(), cr, forge.WithStreams(Streams(t)), forge.WithMountShim()))
}

func TestCloudBuildRunScriptNoShebang(t *testing.T) {
	cr := Runtime(t)

	step := &forge.CloudBuild{
		Step: cloudbuild.Step{
			Name:   "public.ecr.aws/docker/library/alpine",
			Script: `[ "$1" = "hello" ]`,
			Args:   []string{"hello"},
		},
	}

	require.NoError(t, step.Run(t.Context(), cr, forge.WithStreams(Streams(t)), forge.WithMountShim()))
}

func TestCloudBuildRunScriptWithEntrypointFails(t *testing.T) {
	cr := Runtime(t)

	step := &forge.CloudBuild{
		Step: cloudbuild.Step{
			Name:       "public.ecr.aws/docker/library/alpine",
			Script:     "#!/bin/sh\nexit 0\n",
			Entrypoint: "sh",
		},
	}

	require.Error(t, step.Run(t.Context(), cr, forge.WithStreams(Streams(t)), forge.WithMountShim()))
}

func TestCloudBuildRunScriptNode(t *testing.T) {
	cr := Runtime(t)

	step := &forge.CloudBuild{
		Step: cloudbuild.Step{
			Name:   "public.ecr.aws/docker/library/node",
			Script: "#!/usr/bin/env node\nprocess.exit(0);\n",
		},
	}

	require.NoError(t, step.Run(t.Context(), cr, forge.WithStreams(Streams(t)), forge.WithMountShim()))
}

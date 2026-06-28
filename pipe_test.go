//go:build shim && (docker || dockerd)

package forge_test

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/frantjc/forge"
	xos "github.com/frantjc/x/os"
	"github.com/stretchr/testify/require"
)

func TestPipe(t *testing.T) {
	cr := Runtime(t)

	pipe := &forge.Pipe{
		From: &forge.Pure{
			Image:      "public.ecr.aws/docker/library/alpine",
			Entrypoint: []string{"sh", "-c"},
			Cmd:        []string{"exit 0"},
		},
		To: &forge.Pure{
			Image:      "public.ecr.aws/docker/library/alpine",
			Entrypoint: []string{"sh", "-c"},
			Cmd:        []string{"exit 0"},
		},
	}

	require.NoError(t, pipe.Run(t.Context(), cr, forge.WithStreams(Streams(t))))
}

func TestPipeFromNonzeroExitCode(t *testing.T) {
	cr := Runtime(t)

	expected := rand.IntN(254) + 2
	pipe := &forge.Pipe{
		From: &forge.Pure{
			Image:      "public.ecr.aws/docker/library/alpine",
			Entrypoint: []string{"sh", "-c"},
			Cmd:        []string{fmt.Sprintf("exit %d", expected)},
		},
		To: &forge.Pure{
			Image:      "public.ecr.aws/docker/library/alpine",
			Entrypoint: []string{"sh", "-c"},
			Cmd:        []string{"exit 0"},
		},
	}

	err := pipe.Run(t.Context(), cr, forge.WithStreams(Streams(t)))
	require.Error(t, err)
	require.ErrorIs(t, err, forge.ErrContainerExitedWithNonzeroExitCode)
	require.Equal(t, expected, xos.ErrorExitCode(err))
}

func TestPipeToNonzeroExitCode(t *testing.T) {
	cr := Runtime(t)

	expected := rand.IntN(254) + 2
	pipe := &forge.Pipe{
		From: &forge.Pure{
			Image:      "public.ecr.aws/docker/library/alpine",
			Entrypoint: []string{"sh", "-c"},
			Cmd:        []string{"exit 0"},
		},
		To: &forge.Pure{
			Image:      "public.ecr.aws/docker/library/alpine",
			Entrypoint: []string{"sh", "-c"},
			Cmd:        []string{fmt.Sprintf("exit %d", expected)},
		},
	}

	err := pipe.Run(t.Context(), cr, forge.WithStreams(Streams(t)))
	require.Error(t, err)
	require.ErrorIs(t, err, forge.ErrContainerExitedWithNonzeroExitCode)
	require.Equal(t, expected, xos.ErrorExitCode(err))
}

func TestPipeStdout(t *testing.T) {
	cr := Runtime(t)

	notExpected := "hello there"
	expected := "general kenobi"
	require.NotContains(t, expected, notExpected)
	pipe := &forge.Pipe{
		From: &forge.Pure{
			Image:      "public.ecr.aws/docker/library/alpine",
			Entrypoint: []string{"sh", "-c"},
			Cmd: []string{
				fmt.Sprintf("echo %s; echo %s", notExpected, expected),
			},
		},
		To: &forge.Pure{
			Image:      "public.ecr.aws/docker/library/alpine",
			Entrypoint: []string{"sh", "-c"},
			Cmd: []string{
				fmt.Sprintf("grep '%s'", expected),
			},
		},
	}

	ctx := t.Context()
	streams, buf := StreamsCaptureStdout(t)
	require.NoError(t, pipe.Run(ctx, cr, forge.WithStreams(streams)))

	stdout := buf.String()
	require.Equal(t, strings.TrimSpace(stdout), expected)
	require.NotContains(t, stdout, notExpected)
}

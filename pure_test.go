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

func TestPure(t *testing.T) {
	cr := Runtime(t)

	pure := &forge.Pure{
		Image:      "public.ecr.aws/docker/library/alpine",
		Entrypoint: []string{"sh", "-c"},
		Cmd: []string{
			"exit 0",
		},
	}

	require.NoError(t, pure.Run(t.Context(), cr, forge.WithStreams(Streams(t))))
}

func TestPureNonzeroExitCode(t *testing.T) {
	cr := Runtime(t)

	expected := rand.IntN(254) + 2
	pure := &forge.Pure{
		Image:      "public.ecr.aws/docker/library/alpine",
		Entrypoint: []string{"sh", "-c"},
		Cmd: []string{
			fmt.Sprintf("exit %d", expected),
		},
	}
	err := pure.Run(t.Context(), cr, forge.WithStreams(Streams(t)))
	require.Error(t, err)
	require.ErrorIs(t, err, forge.ErrContainerExitedWithNonzeroExitCode)
	require.Equal(t, expected, xos.ErrorExitCode(err))
}

func TestPureStdout(t *testing.T) {
	cr := Runtime(t)

	expected := "hello there general kenobi"
	pure := &forge.Pure{
		Image:      "public.ecr.aws/docker/library/alpine",
		Entrypoint: []string{"sh", "-c"},
		Cmd: []string{
			fmt.Sprintf("echo %s", expected),
		},
	}

	streams, buf := StreamsCaptureStdout(t)
	require.NoError(t, pure.Run(t.Context(), cr, forge.WithStreams(streams)))

	stdout := buf.String()
	require.Equal(t, strings.TrimSpace(stdout), expected)
}

//go:build shim && (docker || dockerd)

package forge_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/runtime"
	"github.com/stretchr/testify/require"
)

func Runtime(t testing.TB) forge.ContainerRuntime {
	t.Helper()
	ctx := t.Context()
	cr, err := runtime.New(ctx, "")
	require.NoError(t, err)
	return cr
}

func DiscardStreams(t testing.TB) *forge.Streams {
	t.Helper()
	return &forge.Streams{Out: io.Discard, Err: io.Discard}
}

func Streams(t testing.TB) *forge.Streams {
	t.Helper()
	return &forge.Streams{Out: t.Output(), Err: t.Output()}
}

func StreamsCaptureStdout(t testing.TB) (*forge.Streams, *bytes.Buffer) {
	t.Helper()
	buf := new(bytes.Buffer)
	return &forge.Streams{Out: buf, Err: t.Output()}, buf
}

func MountShim(t testing.TB) forge.RunOpt {
	t.Helper()
	if os.Getenv("DAGGER_SESSION_TOKEN") == "" {
		return MountShim(t)
	}
	return new(forge.RunOpts)
}

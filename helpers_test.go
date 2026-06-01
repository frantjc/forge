//go:build shim && (docker || dockerd)

package forge_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/runtime"
	"github.com/stretchr/testify/require"
)

func Runtime(t *testing.T) forge.ContainerRuntime {
	t.Helper()
	ctx := t.Context()
	cr, err := runtime.New(ctx, "")
	require.NoError(t, err)
	return cr
}

func DiscardStreams(t *testing.T) *forge.Streams {
	t.Helper()
	return &forge.Streams{Out: io.Discard, Err: io.Discard}
}

func Streams(t *testing.T) *forge.Streams {
	t.Helper()
	return &forge.Streams{Out: t.Output(), Err: t.Output()}
}

func StreamsCaptureStdout(t *testing.T) (*forge.Streams, *bytes.Buffer) {
	t.Helper()
	buf := new(bytes.Buffer)
	return &forge.Streams{Out: buf, Err: t.Output()}, buf
}

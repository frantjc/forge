package forge

import (
	"io"
	"os"
)

var (
	DefaultDetachKeys = "ctrl-d"
)

type Streams struct {
	In         io.Reader
	Out, Err   io.Writer
	Tty        bool
	DetachKeys string
}

func StdStreams() *Streams {
	return &Streams{
		In:         os.Stdin,
		Out:        os.Stdout,
		Err:        os.Stderr,
		// Tty:        true,
		// DetachKeys: DefaultDetachKeys,
	}
}

package forge

import (
	"io"
	"os"
)

type Drains struct {
	Out, Err io.Writer
	Tty      bool
}

func (d *Drains) ToStreams(in io.Reader) *Streams {
	return &Streams{
		In:     in,
		Drains: d,
	}
}

func StdDrains() *Drains {
	return &Drains{
		Out: os.Stdout,
		Err: os.Stderr,
	}
}

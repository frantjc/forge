package forge

import (
	"fmt"
	"io"
	"os"
)

// Drains represents only outward streams from an Ore,
// namely stdout and stderr.
type Drains struct {
	Out, Err io.Writer
	Tty      bool
}

// ToStreams turns a Drains to a Streams for use
// by Ores to pass to a Container.
func (d *Drains) ToStreams(in io.Reader) *Streams {
	return &Streams{
		In:     in,
		Drains: d,
	}
}

// GoString implements fmt.GoStringer.
func (d *Drains) GoString() string {
	return "&Drains{Out: " + fmt.Sprint(d.Out) + ", Err: " + fmt.Sprint(d.Err) + ", Tty: " + fmt.Sprint(d.Tty) + "}"
}

// StdDrains returns a Drains draining to
// os.Stdout and os.Stderr.
func StdDrains() *Drains {
	return &Drains{
		Out: os.Stdout,
		Err: os.Stderr,
	}
}

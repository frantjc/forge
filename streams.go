package forge

import (
	"io"
	"os"

	"github.com/moby/term"
)

var (
	// DefaultDetachKeys are the default key combinations
	// to use when detaching from a Container that has
	// been attached to
	DefaultDetachKeys = "ctrl-d"
)

// Streams represents streams to and from a process
// inside of a Container
type Streams struct {
	*Drains
	In         io.Reader
	DetachKeys string
}

// StdStreams returns a Streams consisting of os.Stdin,
// os.Stdout and os.Stderr
func StdStreams() *Streams {
	return &Streams{
		In:         os.Stdin,
		Drains:     StdDrains(),
		DetachKeys: DefaultDetachKeys,
	}
}

// FileDescriptor is an interface to check io.Readers and io.Writers
// against to inspect if they are terminals
type FileDescriptor interface {
	Fd() uintptr
}

// StdTerminalStreams creates a Streams with os.Stdin, os.Stdout and os.Stderr
// made raw and a restore function to return them to their previous state.
// For use with attaching to a shell inside of a Container
func StdTerminalStreams() (*Streams, func() error) {
	streams, restore, err := TerminalStreams(os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}

	return streams, restore
}

// TerminalStreams creates a Streams with each of the given streams
// that is a terminal made raw and a restore function to return
// them to their previous states. For use with attaching to
// a shell inside of a Container
func TerminalStreams(stdin io.Reader, stdout, stderr io.Writer) (*Streams, func() error, error) {
	var (
		fds    = []FileDescriptor{}
		states = []*term.State{}
	)

	for _, fd := range []interface{}{stdin, stdout, stderr} {
		if fd, ok := fd.(FileDescriptor); ok {
			if term.IsTerminal(fd.Fd()) {
				state, err := term.MakeRaw(fd.Fd())
				if err != nil {
					return nil, nil, err
				}

				states = append(states, state)
				fds = append(fds, fd)
			}

		}
	}

	return &Streams{
			In: stdin,
			Drains: &Drains{
				Out: stdout,
				Err: stderr,
				Tty: true,
			},
			DetachKeys: DefaultDetachKeys,
		}, func() error {
			for i, fd := range fds {
				if err := term.RestoreTerminal(fd.Fd(), states[i]); err != nil {
					return err
				}
			}

			return nil
		}, nil
}

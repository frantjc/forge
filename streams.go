package forge

import (
	"errors"
	"io"
	"os"

	"github.com/moby/term"
)

var (
	DefaultDetachKeys = "ctrl-d"
)

type Streams struct {
	*Drains
	In         io.Reader
	DetachKeys string
}

func StdStreams() *Streams {
	return &Streams{
		In:         os.Stdin,
		Drains:     StdDrains(),
		DetachKeys: DefaultDetachKeys,
	}
}

var (
	ErrNotATerminal = errors.New("not a terminal")
)

type FileDescriptor interface {
	Fd() uintptr
}

func StdTerminalStreams() (*Streams, func() error) {
	streams, restore, err := TerminalStreams(os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}

	return streams, restore
}

func TerminalStreams(stdin io.Reader, stdout, stderr io.Writer) (*Streams, func() error, error) {
	var (
		fds    = []FileDescriptor{}
		states = []*term.State{}
	)

	for _, fd := range []interface{}{stdin, stdout, stderr} {
		if fd, ok := fd.(FileDescriptor); ok {
			if !term.IsTerminal(fd.Fd()) {
				return nil, nil, ErrNotATerminal
			}

			state, err := term.MakeRaw(fd.Fd())
			if err != nil {
				return nil, nil, err
			}

			states = append(states, state)
			fds = append(fds, fd)
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

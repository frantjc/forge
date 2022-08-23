package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/frantjc/forge"
	"github.com/moby/term"
)

func (c *Container) Run(ctx context.Context, streams *forge.Streams) (int, error) {
	var (
		stdin          io.Reader
		stdout, stderr io.Writer
		detachKeys     string
		tty            bool
	)
	if streams != nil {
		stdin = streams.In
		stdout = streams.Out
		stderr = streams.Err
		tty = streams.Tty
		if tty {
			stderr = stdout
		}
		detachKeys = streams.DetachKeys
	}


	hjr, err := c.ContainerAttach(ctx, c.ID, types.ContainerAttachOptions{
		Stream:     streams != nil,
		Stdin:      stdin != nil,
		Stdout:     stdout != nil,
		Stderr:     stderr != nil,
		DetachKeys: detachKeys,
	})
	if err != nil {
		return -1, err
	}

	errC := make(chan error, 1)
	go func() {
		var err error
		if tty {
			_, err = io.Copy(stdout, hjr.Reader)
		} else {
			_, err = stdcopy.StdCopy(
				stdout,
				stderr,
				hjr.Reader,
			)
		}
		errC <- err
	}()

	if stdin != nil {
		if detachKeys != "" {
			detachKeysB, err := term.ToBytes(detachKeys)
			if err != nil {
				return -1, err
			}

			stdin = term.NewEscapeProxy(stdin, detachKeysB)
		}

		go func() {
			if _, err = io.Copy(hjr.Conn, stdin); err != nil {
				errC <- err
			}

			errC <- hjr.CloseWrite()
		}()
	}

	if err = c.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
		return -1, err
	}

	cwokbC, waitErrC := c.ContainerWait(ctx, c.ID, container.WaitConditionNotRunning)
	if err != nil {
		return -1, err
	}

	select {
	case cwokb := <-cwokbC:
		if cwokb.Error != nil {
			err = fmt.Errorf(cwokb.Error.Message)
		}

		return int(cwokb.StatusCode), err
	case err = <-errC:
		return -1, err
	case err = <-waitErrC:
		return -1, err
	case <-ctx.Done():
		return -1, ctx.Err()
	}
}

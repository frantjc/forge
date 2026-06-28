package dockerd

import (
	"context"
	"io"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/frantjc/forge"
	"github.com/moby/term"
	"golang.org/x/sync/errgroup"
)

type Container struct {
	ID string
	*client.Client
}

func (c *Container) GetID() string {
	return c.ID
}

func (c *Container) GoString() string {
	return "&Container{" + c.GetID() + "}"
}

func (c *Container) CopyTo(ctx context.Context, destination string, content io.Reader) error {
	if rc, ok := content.(io.ReadCloser); ok {
		defer rc.Close()
	}

	return c.CopyToContainer(ctx, c.ID, filepath.Clean(destination), content, container.CopyToContainerOptions{})
}

func (c *Container) CopyFrom(ctx context.Context, source string) (io.ReadCloser, error) {
	rc, _, err := c.CopyFromContainer(ctx, c.ID, filepath.Clean(source))
	return rc, err
}

func (c *Container) Start(ctx context.Context) error {
	return c.ContainerStart(ctx, c.ID, container.StartOptions{})
}

func (c *Container) Exec(ctx context.Context, containerConfig *forge.ContainerConfig, streams *forge.Streams) (int, error) {
	var (
		stdin          io.Reader
		stdout, stderr io.Writer
		tty            bool
		detachKeys     string
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

	idr, err := c.ContainerExecCreate(ctx, c.ID, container.ExecOptions{
		User:         containerConfig.User,
		Privileged:   containerConfig.Privileged,
		Env:          containerConfig.Env,
		WorkingDir:   containerConfig.WorkingDir,
		Cmd:          append(containerConfig.Entrypoint, containerConfig.Cmd...),
		Tty:          tty,
		DetachKeys:   detachKeys,
		AttachStdin:  stdin != nil,
		AttachStdout: stdout != nil,
		AttachStderr: stderr != nil,
	})
	if err != nil {
		return -1, err
	}

	hjr, err := c.ContainerExecAttach(ctx, idr.ID, container.ExecStartOptions{
		Tty: tty,
	})
	if err != nil {
		return -1, err
	}
	defer hjr.Close()

	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if tty {
			if _, err := io.Copy(stdout, hjr.Reader); err != nil {
				return err
			}

			return nil
		}

		if _, err := stdcopy.StdCopy(stdout, stderr, hjr.Reader); err != nil {
			return err
		}

		return nil
	})

	if stdin != nil {
		eg.Go(func() error {
			defer hjr.CloseWrite()

			_stdin := stdin

			if detachKeys != "" {
				detachKeysB, err := term.ToBytes(detachKeys)
				if err != nil {
					return err
				}
				_stdin = term.NewEscapeProxy(_stdin, detachKeysB)
			}

			if _, err := io.Copy(hjr.Conn, _stdin); err != nil {
				return err
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return -1, err
	}

	cei, err := c.ContainerExecInspect(ctx, idr.ID)
	if err != nil {
		return -1, err
	}

	return cei.ExitCode, nil
}

func (c *Container) Stop(ctx context.Context) error {
	seconds := 0
	if deadline, ok := ctx.Deadline(); ok {
		seconds = int(time.Until(deadline).Seconds())
	}

	return c.ContainerStop(ctx, c.ID, container.StopOptions{
		Timeout: &seconds,
	})
}

func (c *Container) Remove(ctx context.Context) error {
	return c.ContainerRemove(ctx, c.ID, container.RemoveOptions{
		Force: true,
	})
}

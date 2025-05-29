package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/changroup"
	"github.com/moby/term"
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

	hjr, err := c.ContainerAttach(ctx, c.ID, container.AttachOptions{
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
		if tty {
			_, err = io.Copy(stdout, hjr.Reader)
		} else {
			_, err = stdcopy.StdCopy(
				stdout,
				stderr,
				hjr.Reader,
			)
		}
		if err != nil {
			errC <- err
		}
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

			if err = hjr.CloseWrite(); err != nil {
				errC <- hjr.CloseWrite()
			}
		}()
	}

	if err = c.ContainerStart(ctx, c.ID, container.StartOptions{}); err != nil {
		return -1, err
	}

	cwokbC, waitErrC := c.ContainerWait(ctx, c.ID, container.WaitConditionNotRunning)

	select {
	case cwokb := <-cwokbC:
		if cwokb.Error != nil {
			err = fmt.Errorf("%s", cwokb.Error.Message)
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

	errC := make(chan error, 1)
	outC := make(chan any, 1)
	go func() {
		if tty {
			_, err = io.Copy(stdout, hjr.Reader)
		} else {
			// "exit status 1" comes from here
			_, err = stdcopy.StdCopy(
				stdout,
				stderr,
				hjr.Reader,
			)
		}
		if err != nil {
			errC <- err
		}

		go close(outC)
	}()

	inC := make(chan any, 1)
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

			if err = hjr.CloseWrite(); err != nil {
				errC <- err
			}

			go close(inC)
		}()
	} else {
		go close(inC)
	}

	select {
	case err = <-errC:
		if _, ok := err.(term.EscapeError); ok {
			err = nil
		}
	case <-ctx.Done():
		err = ctx.Err()
	case <-changroup.AllSettled(inC, outC):
	}
	if err != nil {
		return -1, err
	}

	cei, inspectErr := c.ContainerExecInspect(ctx, idr.ID)
	if inspectErr != nil {
		return -1, inspectErr
	}

	return cei.ExitCode, err
}

func (c *Container) Restart(ctx context.Context) error {
	seconds := -1
	if deadline, ok := ctx.Deadline(); ok {
		seconds = int(time.Until(deadline).Seconds())
	}

	return c.ContainerRestart(ctx, c.ID, container.StopOptions{
		Timeout: &seconds,
	})
}

func (c *Container) Stop(ctx context.Context) error {
	seconds := -1
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

func (c *Container) Kill(ctx context.Context) error {
	return c.ContainerKill(ctx, c.ID, os.Kill.String())
}

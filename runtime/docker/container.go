package docker

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/frantjc/forge"
	xos "github.com/frantjc/x/os"
)

type Container struct {
	ID   string
	Path string
}

func (c *Container) GetID() string { return c.ID }

func (c *Container) CopyTo(ctx context.Context, dst string, r io.Reader) error {
	//nolint:gosec
	cmd := exec.CommandContext(ctx, c.Path, "cp", "-", fmt.Sprintf("%s:%s", c.ID, dst))
	cmd.Stdin = r

	if err := cmd.Run(); err != nil {
		return xos.NewExitCodeError(err, cmd.ProcessState.ExitCode())
	}

	return nil
}

func (c *Container) CopyFrom(ctx context.Context, src string) (io.ReadCloser, error) {
	//nolint:gosec
	cmd := exec.CommandContext(ctx, c.Path, "cp", fmt.Sprintf("%s:%s", c.ID, src), "-")

	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Container) Start(ctx context.Context) error {
	//nolint:gosec
	cmd := exec.CommandContext(ctx, c.Path, "start", c.ID)

	if err := cmd.Run(); err != nil {
		return xos.NewExitCodeError(err, cmd.ProcessState.ExitCode())
	}

	return nil
}

func (c *Container) Exec(ctx context.Context, cfg *forge.ContainerConfig, streams *forge.Streams) (int, error) {
	args := []string{"exec"}

	if cfg.User != "" {
		args = append(args, "--user", cfg.User)
	}

	if cfg.WorkingDir != "" {
		args = append(args, "--workdir", cfg.WorkingDir)
	}

	for _, env := range cfg.Env {
		args = append(args, "--env", env)
	}

	args = append(args, c.ID)
	args = append(args, cfg.Entrypoint...)
	args = append(args, cfg.Cmd...)

	//nolint:gosec
	cmd := exec.CommandContext(ctx, c.Path, args...)
	cmd.Stdin = streams.In
	cmd.Stdout = streams.Out
	cmd.Stderr = streams.Err

	if err := cmd.Run(); err != nil {
		return -1, xos.NewExitCodeError(err, cmd.ProcessState.ExitCode())
	}

	return cmd.ProcessState.ExitCode(), nil
}

func (c *Container) Stop(ctx context.Context) error {
	//nolint:gosec
	return exec.CommandContext(ctx, c.Path, "stop", c.ID).Run()
}

func (c *Container) Remove(ctx context.Context) error {
	//nolint:gosec
	return exec.CommandContext(ctx, c.Path, "rm", "--force", c.ID).Run()
}

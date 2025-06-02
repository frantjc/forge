package docker

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/frantjc/forge"
	xos "github.com/frantjc/x/os"
)

type DockerContainer struct {
	ID   string
	Path string
}

func (c *DockerContainer) GetID() string { return c.ID }

func (c *DockerContainer) CopyTo(ctx context.Context, dst string, r io.Reader) error {
	cmd := exec.CommandContext(ctx, c.Path, "cp", "-", fmt.Sprintf("%s:%s", c.ID, dst))
	cmd.Stdin = r

	if err := cmd.Run(); err != nil {
		return xos.NewExitCodeError(err, cmd.ProcessState.ExitCode())
	}

	return nil
}

func (c *DockerContainer) CopyFrom(ctx context.Context, src string) (io.ReadCloser, error) {
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

func (c *DockerContainer) Start(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, c.Path, "start", c.ID)

	if err := cmd.Run(); err != nil {
		return xos.NewExitCodeError(err, cmd.ProcessState.ExitCode())
	}

	return nil
}

func (c *DockerContainer) Exec(ctx context.Context, cfg *forge.ContainerConfig, streams *forge.Streams) (int, error) {
	args := []string{"exec"}

	if cfg.User != "" {
		args = append(args, "-u", cfg.User)
	}

	if cfg.WorkingDir != "" {
		args = append(args, "-w", cfg.WorkingDir)
	}

	for _, env := range cfg.Env {
		args = append(args, "-e", env)
	}

	args = append(args, c.ID)
	args = append(args, cfg.Entrypoint...)
	args = append(args, cfg.Cmd...)

	cmd := exec.CommandContext(ctx, c.Path, args...)
	cmd.Stdin = streams.In
	cmd.Stdout = streams.Out
	cmd.Stderr = streams.Err

	if err := cmd.Run(); err != nil {
		return -1, xos.NewExitCodeError(err, cmd.ProcessState.ExitCode())
	}

	return cmd.ProcessState.ExitCode(), nil
}

func (c *DockerContainer) Stop(ctx context.Context) error {
	return exec.CommandContext(ctx, c.Path, "stop", c.ID).Run()
}

func (c *DockerContainer) Remove(ctx context.Context) error {
	return exec.CommandContext(ctx, c.Path, "rm", "-f", c.ID).Run()
}

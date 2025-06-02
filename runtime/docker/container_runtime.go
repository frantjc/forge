package docker

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/frantjc/forge"
	xos "github.com/frantjc/x/os"
	"github.com/google/go-containerregistry/pkg/name"
)

type ContainerRuntime struct {
	Path string
}

func New(path string) *ContainerRuntime {
	r := &ContainerRuntime{
		Path: path,
	}

	if r.Path == "" {
		r.Path = "docker"
	}

	return r
}

func (r *ContainerRuntime) GetContainer(ctx context.Context, id string) (forge.Container, error) {
	return &DockerContainer{ID: id, Path: r.Path}, nil
}

func (r *ContainerRuntime) CreateContainer(ctx context.Context, img forge.Image, cfg *forge.ContainerConfig) (forge.Container, error) {
	args := []string{"create"}

	for _, env := range cfg.Env {
		args = append(args, "--env", env)
	}

	for _, m := range cfg.Mounts {
		args = append(args, "--volume", fmt.Sprintf("%s:%s", m.Source, m.Destination))
	}

	if cfg.WorkingDir != "" {
		args = append(args, "--workdir", cfg.WorkingDir)
	}

	if cfg.User != "" {
		args = append(args, "--user", cfg.User)
	}

	if cfg.Privileged {
		args = append(args, "--privileged")
	}

	args = append(args, img.Name())
	args = append(args, cfg.Entrypoint...)
	args = append(args, cfg.Cmd...)

	cmd := exec.CommandContext(ctx, r.Path, args...)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, xos.NewExitCodeError(err, cmd.ProcessState.ExitCode())
	}

	id := strings.TrimSpace(buf.String())

	return &DockerContainer{ID: id, Path: r.Path}, nil
}

func (r *ContainerRuntime) PullImage(ctx context.Context, reference string) (forge.Image, error) {
	ref, err := name.ParseReference(reference)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, r.Path, "pull", ref.String())

	if err := cmd.Run(); err != nil {
		return nil, xos.NewExitCodeError(err, cmd.ProcessState.ExitCode())
	}

	return &Image{Ref: ref.String(), Path: r.Path}, nil
}

func (r *ContainerRuntime) Close() error {
	return nil
}

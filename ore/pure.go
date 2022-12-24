package ore

import (
	"bytes"
	"context"

	"github.com/frantjc/forge"
	cfs "github.com/frantjc/forge/internal/containerfs"
	"github.com/frantjc/forge/internal/containerutil"
	"github.com/frantjc/forge/internal/contaminate"
)

// Pure is an Ore for running a "pure" command inside
// of a container.
type Pure struct {
	Image      string   `json:"image,omitempty"`
	Entrypoint []string `json:"entrypoint,omitempty"`
	Cmd        []string `json:"cmd,omitempty"`
	Env        []string `json:"env,omitempty"`
	Input      []byte   `json:"input,omitempty"`
}

func (o *Pure) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) (*forge.Metal, error) {
	image, err := containerRuntime.PullImage(ctx, o.Image)
	if err != nil {
		return nil, err
	}

	containerConfig := &forge.ContainerConfig{
		Entrypoint: o.Entrypoint,
		Cmd:        o.Cmd,
		Env:        o.Env,
		WorkingDir: cfs.WorkingDir,
		Mounts:     contaminate.MountsFrom(ctx),
	}

	container, err := containerutil.CreateSleepingContainer(ctx, containerRuntime, image, containerConfig)
	if err != nil {
		return nil, err
	}
	defer container.Stop(ctx)   //nolint:errcheck
	defer container.Remove(ctx) //nolint:errcheck

	input := contaminate.InputFrom(ctx)
	if len(input) == 0 {
		input = o.Input
	}

	exitCode, err := container.Exec(ctx, containerConfig, drains.ToStreams(bytes.NewReader(input)))
	if err != nil {
		return nil, err
	}

	return &forge.Metal{
		ExitCode: int64(exitCode),
	}, nil
}

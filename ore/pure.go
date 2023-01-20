package ore

import (
	"bytes"
	"context"

	"github.com/frantjc/forge"
	cfs "github.com/frantjc/forge/internal/containerfs"
	"github.com/frantjc/forge/internal/containerutil"
	"github.com/frantjc/forge/internal/contaminate"
	errorcode "github.com/frantjc/go-error-code"
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

func (o *Pure) Liquify(ctx context.Context, containerRuntime forge.ContainerRuntime, drains *forge.Drains) error {
	image, err := containerRuntime.PullImage(ctx, o.Image)
	if err != nil {
		return err
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
		return err
	}
	defer container.Stop(ctx)   //nolint:errcheck
	defer container.Remove(ctx) //nolint:errcheck

	input := contaminate.InputFrom(ctx)
	if len(input) == 0 {
		input = o.Input
	}

	if exitCode, err := container.Exec(ctx, containerConfig, drains.ToStreams(bytes.NewReader(input))); err != nil {
		return err
	} else if exitCode > 0 {
		return errorcode.New(ErrContainerExitedWithNonzeroExitCode, errorcode.WithExitCode(exitCode))
	}

	return nil
}

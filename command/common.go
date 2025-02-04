package command

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/runtime/docker"
	xslice "github.com/frantjc/x/slice"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func commandStreams(cmd *cobra.Command, stdoutUsed ...bool) *forge.Streams {
	return &forge.Streams{
		In: cmd.InOrStdin(),
		Out: func() io.Writer {
			if xslice.Some(stdoutUsed, func(b bool, _ int) bool {
				return b
			}) {
				return cmd.ErrOrStderr()
			}

			return cmd.OutOrStdout()
		}(),
		Err: cmd.ErrOrStderr(),
	}
}

func hookAttach(cmd *cobra.Command, workingDir string, stdoutUsed ...bool) func(context.Context, forge.Container) {
	return func(ctx context.Context, c forge.Container) {
		var (
			streams = commandStreams(cmd, stdoutUsed...)
			_, _    = fmt.Fprintln(streams.Out, "detach with", forge.DefaultDetachKeys)
		)

		streams, restore, err := forge.TerminalStreams(streams.In, streams.Out, streams.Err)
		defer restore() //nolint:errcheck
		if err != nil {
			return
		}

		for _, shell := range []string{"bash", "sh"} {
			if _, err = c.Exec(
				ctx,
				&forge.ContainerConfig{
					Entrypoint: []string{shell},
					WorkingDir: workingDir,
				},
				streams,
			); err == nil {
				break
			}
		}

		_, _ = fmt.Fprintln(streams.Out)
	}
}

func oreOptsAndContainerRuntime(cmd *cobra.Command) (forge.ContainerRuntime, *forge.OreOpts, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, nil, err
	}

	var (
		ctrWorkDir = "/" + uuid.NewString()
		dindPath   = ctrWorkDir
	)
	if cmd.Flag("no-dind").Changed {
		dindPath = ""
	}

	return docker.New(cli, dindPath), &forge.OreOpts{
		Streams:             commandStreams(cmd),
		InterceptDockerSock: cmd.Flag("fix-dind").Changed,
		WorkingDir:          ctrWorkDir,
	}, nil
}

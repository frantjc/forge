package command

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/client"
	"github.com/frantjc/forge"
	"github.com/frantjc/forge/runtime/dockerd"
	xslice "github.com/frantjc/x/slice"
	"github.com/spf13/cobra"
)

func commandStreams(cmd *cobra.Command, stdoutUsed ...bool) *forge.Streams {
	return &forge.Streams{
		Out: func() io.Writer {
			if xslice.Some(stdoutUsed, func(b bool, _ int) bool {
				return b
			}) {
				return cmd.ErrOrStderr()
			}

			return cmd.OutOrStdout()
		}(),
		Err:        cmd.ErrOrStderr(),
		DetachKeys: forge.DefaultDetachKeys,
	}
}

func hookAttach(cmd *cobra.Command, workingDir string, stdoutUsed ...bool) func(context.Context, forge.Container) {
	return func(ctx context.Context, c forge.Container) {
		var (
			streams = commandStreams(cmd, stdoutUsed...)
			_, _    = fmt.Fprintln(streams.Out, "detach with", streams.DetachKeys)
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

func runOptsAndContainerRuntime(cmd *cobra.Command, stdoutUsed ...bool) (forge.ContainerRuntime, *forge.RunOpts, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, nil, err
	}

	var (
		ctrWorkDir = "/forge"
		dindPath   = ctrWorkDir
	)
	if cmd.Flag("no-dind").Changed {
		dindPath = ""
	}

	return dockerd.New(cli, dindPath),
		&forge.RunOpts{
			Streams:             commandStreams(cmd, stdoutUsed...),
			InterceptDockerSock: cmd.Flag("fix-dind").Changed,
			WorkingDir:          ctrWorkDir,
		}, nil
}

func setCommon(cmd *cobra.Command) *cobra.Command {
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.Flags().BoolP("help", "h", false, "Help for "+cmd.Name())
	cmd.Flags().Bool("fix-dind", false, "Intercept and fix traffic to docker.sock")
	cmd.Flags().Bool("no-dind", false, "Disable Docker in Docker")
	cmd.MarkFlagsMutuallyExclusive("no-dind", "fix-dind")

	return cmd
}

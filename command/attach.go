package command

import (
	"context"
	"fmt"

	"github.com/frantjc/forge"
	"github.com/spf13/cobra"
)

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

func commandStreams(cmd *cobra.Command, stdoutUsed ...bool) *forge.Streams {
	return commandDrains(cmd, stdoutUsed...).ToStreams(cmd.InOrStdin())
}

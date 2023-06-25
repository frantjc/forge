package command

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/containerfs"
	"github.com/spf13/cobra"
)

func hookAttach(cmd *cobra.Command, stdoutUsed ...bool) func(context.Context, forge.Container) {
	return func(ctx context.Context, c forge.Container) {
		var (
			streams = commandStreams(cmd, stdoutUsed...)
			_, _    = streams.Out.Write([]byte("detach with " + forge.DefaultDetachKeys + "\n"))
		)

		streams, restore, err := forge.TerminalStreams(streams.In, streams.Out, streams.Err)
		if err != nil {
			return
		}
		defer restore() //nolint:errcheck

		_, _ = c.Exec(
			ctx,
			&forge.ContainerConfig{
				Entrypoint: []string{"sh"},
				WorkingDir: containerfs.WorkingDir,
			},
			streams,
		)

		_, _ = streams.Out.Write([]byte("\n"))
	}
}

func commandStreams(cmd *cobra.Command, stdoutUsed ...bool) *forge.Streams {
	return commandDrains(cmd, stdoutUsed...).ToStreams(cmd.InOrStdin())
}

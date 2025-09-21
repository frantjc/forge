package command

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/frantjc/forge/githubactions"
	"github.com/spf13/cobra"
)

// NewShim returns the command which acts as
// the entrypoint for `shim`.
func NewShim() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "shim",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			wd, err := os.Getwd()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}

			globalContext, err := githubactions.NewGlobalContextFromPath(wd)
			if err != nil {
				globalContext = githubactions.NewGlobalContextFromEnv()
			}

			subcmd := exec.CommandContext(ctx, args[0], args[1:]...) //nolint:gosec
			subcmd.Dir = globalContext.GitHubContext.Workspace
			subcmd.Env = append(os.Environ(), globalContext.Env()...)
			subcmd.Stdin = cmd.InOrStdin()
			subcmd.Stdout = githubactions.NewWorkflowCommandWriter(cmd.OutOrStdout(), globalContext)
			subcmd.Stderr = cmd.ErrOrStderr()

			return subcmd.Run()
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Help for "+cmd.Name())
	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")

	return cmd
}

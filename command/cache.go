package command

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/frantjc/forge/internal/hostfs"
	xslices "github.com/frantjc/x/slices"
	"github.com/spf13/cobra"
)

var (
	runnerTmpArgs           = []string{"runner_temp", "runnertemp", "runner_tmp", "runnertmp", "temp", "tmp"}
	runnerToolCacheArgs     = []string{"runner_tool_cache", "runner_toolcache", "runnertoolcache", "toolcache", "tc"}
	actionsArgs             = []string{"github", "github_actions", "githubactions", "gha", "action", "actions"}
	cloudbuildWorkspaceArgs = []string{"cloudbuild", "cb", "workspace"}
)

// NewCache returns the command which acts as
// the entrypoint for `forge cache`.
func NewCache() *cobra.Command {
	var (
		clean bool
		cmd   = setCommon(&cobra.Command{
			Use:       "cache [name] [--clean]",
			Short:     "Interact with the Forge cache",
			Args:      cobra.MaximumNArgs(1),
			ValidArgs: append(runnerTmpArgs, runnerToolCacheArgs...),
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					cache = hostfs.CacheHome
					arg   = strings.ToLower(xslices.Find(args, func(_ string, _ int) bool {
						return true
					}))
				)

				switch {
				case arg == "":
				case slices.Contains(runnerTmpArgs, arg):
					cache = hostfs.RunnerTmp
				case slices.Contains(runnerToolCacheArgs, arg):
					cache = hostfs.RunnerToolCache
				case slices.Contains(actionsArgs, arg):
					cache = hostfs.ActionsCache
				case slices.Contains(cloudbuildWorkspaceArgs, arg):
					cache = hostfs.CloudBuildWorkspace
				default:
					return fmt.Errorf("unknown cache: %s", arg)
				}

				if clean {
					return os.RemoveAll(cache)
				}

				if err := os.MkdirAll(cache, 0o755); err != nil {
					return err
				}

				_, err := fmt.Fprintln(cmd.OutOrStdout(), cache)
				return err
			},
		})
	)

	cmd.Flags().BoolVar(&clean, "clean", false, "Clean the cache and exit")

	return cmd
}

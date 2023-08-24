package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/frantjc/forge/internal/hostfs"
	"github.com/frantjc/go-fn"
	"github.com/spf13/cobra"
)

var (
	runnerTmpArgs       = []string{"runner_temp", "runnertemp", "runner_tmp", "runnertmp", "temp", "tmp"}
	runnerToolCacheArgs = []string{"runner_tool_cache", "runner_toolcache", "runnertoolcache", "toolcache"}
	actionsArgs         = []string{"github", "github_actions", "githubactions", "gha", "action", "actions"}
	oreArgs             = []string{"ore", "ores"}
)

// NewCache returns the command which acts as
// the entrypoint for `forge cache`.
func NewCache() *cobra.Command {
	var (
		clean bool
		cmd   = &cobra.Command{
			Use:           "cache",
			Short:         "Interact with the Forge cache",
			SilenceErrors: true,
			SilenceUsage:  true,
			Args:          cobra.MaximumNArgs(1),
			ValidArgs:     append(runnerTmpArgs, runnerToolCacheArgs...),
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					cache = hostfs.CacheHome
					arg   = strings.ToLower(fn.Find(args, func(_ string, _ int) bool {
						return true
					}))
				)

				switch {
				case arg == "":
				case fn.Includes(runnerTmpArgs, arg):
					cache = hostfs.RunnerTmp
				case fn.Includes(runnerToolCacheArgs, arg):
					cache = hostfs.RunnerToolCache
				case fn.Includes(actionsArgs, arg):
					cache = hostfs.ActionsCache
				case fn.Includes(oreArgs, arg):
					cache = hostfs.OreCache
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
		}
	)

	cmd.PersistentFlags().BoolVar(&clean, "clean", false, "clean the cache")

	return cmd
}

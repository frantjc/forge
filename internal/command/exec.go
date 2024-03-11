package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/frantjc/forge/envconv"
	"github.com/frantjc/forge/githubactions"
	xos "github.com/frantjc/x/os"
	"github.com/spf13/cobra"
)

// NewExec returns the command which acts as
// the entrypoint for `shim exec`.
func NewExec() *cobra.Command {
	var (
		workingdir string
		cmd        = &cobra.Command{
			Use:           "exec",
			Short:         "Execute the given command after sourcing $GITHUB_PATH and $GITHUB_ENV, if set",
			SilenceErrors: true,
			SilenceUsage:  true,
			Args:          cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx              = cmd.Context()
					command          = exec.CommandContext(ctx, args[0], args[1:]...) //nolint:gosec
					githubEnvPath    = os.Getenv(githubactions.EnvVarEnv)
					githubPathPath   = os.Getenv(githubactions.EnvVarPath)
					githubStatePath  = os.Getenv(githubactions.EnvVarState)
					githubOutputPath = os.Getenv(githubactions.EnvVarOutput)
				)

				command.Env = os.Environ()
				command.Stdin = os.Stdin
				command.Stdout = os.Stdout
				command.Stderr = os.Stderr

				for _, githubPath := range []string{
					githubStatePath,
					githubOutputPath,
					githubEnvPath,
					githubPathPath,
				} {
					if githubPath != "" {
						if err := os.MkdirAll(filepath.Dir(githubPath), 0o755); err != nil {
							return err
						}

						if _, err := os.Create(githubPath); err != nil {
							return err
						}
					}
				}

				if githubEnvPath != "" {
					if file, err := os.Open(githubEnvPath); err == nil {
						if githubEnv, err := githubactions.ParseEnvFile(file); err == nil {
							command.Env = append(command.Env, envconv.MapToArr(githubEnv)...)
						}
					}
				}

				path := "PATH=" + os.Getenv("PATH")
				if runnerToolCache := os.Getenv(githubactions.EnvVarRunnerToolCache); runnerToolCache != "" {
					path += ":" + runnerToolCache
				}

				if githubPathPath != "" {
					if file, err := os.Open(githubPathPath); err == nil {
						if githubPath, err := githubactions.ParsePathFile(file); err == nil {
							path += ":" + githubPath
						}
					}
				}

				forgeSock := filepath.Join(workingdir, "forge.sock")
				_, err := os.Stat(forgeSock)
				useForgeSock := workingdir != "" && err == nil
				dockerHost := fmt.Sprintf("DOCKER_HOST=unix://%s", forgeSock)

				var (
					injectedPath       = false
					injectedDockerHost = false
				)

				for i, env := range command.Env {
					parts := strings.SplitN(env, "=", 2)
					if len(parts) > 0 {
						key := parts[0]
						if !injectedPath && strings.EqualFold(key, "PATH") {
							command.Env[i] = path
							injectedPath = true
						} else if useForgeSock && !injectedDockerHost && strings.EqualFold(key, "DOCKER_HOST") {
							command.Env[i] = dockerHost
							injectedDockerHost = true
						}
					}
				}

				if !injectedPath {
					command.Env = append(command.Env, path)
				}

				if useForgeSock && !injectedDockerHost {
					command.Env = append(command.Env, dockerHost)
				}

				// TODO: Wait on forge.sock to be ready.

				return xos.NewExitCodeError(command.Run(), command.ProcessState.ExitCode())
			},
		}
	)

	cmd.Flags().StringVar(&workingdir, "wd", "", "working directory for forge")
	_ = cmd.MarkFlagDirname("wd")

	return cmd
}


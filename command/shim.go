package command

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/envconv"
	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/internal/dind"
	xos "github.com/frantjc/x/os"
	"github.com/spf13/cobra"
)

// NewShim returns the command which acts as
// the entrypoint for `shim`.
func NewShim() *cobra.Command {
	var (
		verbosity int
		cmd       = &cobra.Command{
			Use:           "shim",
			Version:       forge.SemVer(),
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRun: func(cmd *cobra.Command, _ []string) {
				cmd.SetContext(
					forge.WithLogger(cmd.Context(), forge.NewLogger().V(2-verbosity)),
				)
			},
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "V", "verbosity for forge")
	cmd.AddCommand(NewSleep(), NewExec())

	return cmd
}

// NewSleep returns the command which acts as
// the entrypoint for `shim sleep`.
func NewSleep() *cobra.Command {
	var (
		workingdir string
		mounts     map[string]string
		cmd        = &cobra.Command{
			Use:           "sleep",
			Short:         "Sleep until signalled",
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, _ []string) error {
				ctx := cmd.Context()

				if len(mounts) > 0 && workingdir != "" {
					forgeSock := filepath.Join(workingdir, "forge.sock")

					if lis, err := net.Listen("unix", forgeSock); err == nil {
						defer os.Remove(forgeSock)

						if dockerHost := os.Getenv("DOCKER_HOST"); dockerHost != "" {
							if dockerSock, err := url.Parse(dockerHost); err == nil {
								return dind.NewProxy(ctx, mounts, lis, dockerSock)
							}
						}
					}
				}

				<-ctx.Done()

				return ctx.Err()
			},
		}
	)

	cmd.Flags().StringToStringVar(&mounts, "mount", nil, "mounts for forge")
	cmd.Flags().StringVar(&workingdir, "wd", "", "working directory for forge")
	_ = cmd.MarkFlagDirname("wd")

	return cmd
}

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

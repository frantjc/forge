package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/frantjc/forge/envconv"
	"github.com/frantjc/forge/githubactions"
	xos "github.com/frantjc/x/os"
	xslice "github.com/frantjc/x/slice"
)

var (
	errHelp = errors.New("help")
	help    = fmt.Sprintf(`
%s [-s|-e|-h] [args]

  -s   sleep
  -e   execute the given command after sourcing $GITHUB_PATH and $GITHUB_ENV, if set
  -h   help

`, os.Args[0])
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	if err := mainE(ctx); errors.Is(err, errHelp) {
		os.Stderr.WriteString(help)
	} else {
		stop()
		xos.ExitFromError(err)
	}
}

func mainE(ctx context.Context) error {
	args := os.Args

	if len(args) < 2 {
		return errHelp
	}

	switch args[1] {
	// sleep
	case "-s":
		_, _ = os.Stdout.WriteString("zzz...\n")
		<-ctx.Done()
		return ctx.Err()
	// exec
	case "-e":
		if len(args) < 3 {
			return errHelp
		}

		var (
			command          = exec.CommandContext(ctx, args[2], args[3:]...) //nolint:gosec
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

		if i := xslice.FindIndex(command.Env, func(s string, _ int) bool {
			spl := strings.Split(s, "=")
			return len(spl) > 0 && strings.EqualFold(spl[0], "PATH")
		}); i >= 0 {
			command.Env[i] = path
		} else {
			command.Env = append(command.Env, path)
		}

		return xos.NewExitCodeError(command.Run(), command.ProcessState.ExitCode())
	}

	return errHelp
}

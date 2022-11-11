package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/frantjc/forge/pkg/envconv"
	"github.com/frantjc/forge/pkg/errbubble"
	"github.com/frantjc/forge/pkg/githubactions"
	"github.com/frantjc/go-fn"
)

var (
	errHelp = errors.New("help")
	help    = fmt.Sprintf(`
%s [-s|-e|-h] [args]

  -s   sleep
  -e   execute the given command after sourcing $GITHUB_PATH and $GITHUB_ENV
  -h   help

`, os.Args[0])
)

func main() {
	var (
		ctx, stop = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		err       error
	)

	if err := mainE(ctx); err != nil {
		if errors.Is(err, errHelp) {
			os.Stderr.WriteString(help)
			err = nil
		} else {
			os.Stderr.WriteString(err.Error() + "\n")
		}
	}

	stop()
	os.Exit(errbubble.ExitCode(err))
}

func mainE(ctx context.Context) error {
	var (
		args = os.Args
	)

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
			command        = exec.CommandContext(ctx, args[2], args[3:]...) //nolint:gosec
			githubEnvPath  = os.Getenv(githubactions.EnvVarEnv)
			githubPathPath = os.Getenv(githubactions.EnvVarPath)
		)

		command.Env = os.Environ()
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if githubEnv, err := envconv.ArrFromFile(githubEnvPath); err == nil {
			command.Env = append(command.Env, githubEnv...)
		} else {
			if _, err = os.Create(githubEnvPath); err != nil {
				return err
			}
		}

		var (
			path = "PATH=" + os.Getenv("PATH")
		)
		if runnerToolCache := os.Getenv(githubactions.EnvVarRunnerToolCache); runnerToolCache != "" {
			path += ":" + runnerToolCache
		}

		if githubPath, err := envconv.PathFromFile(githubPathPath); err == nil && githubPath != "" {
			path += ":" + githubPath
		} else {
			if _, err = os.Create(githubPathPath); err != nil {
				return err
			}
		}

		if i := fn.FindIndex(command.Env, func(s string, _ int) bool {
			spl := strings.Split(s, "=")
			return len(spl) > 0 && strings.EqualFold(spl[0], "PATH")
		}); i >= 0 {
			command.Env[i] = path
		} else {
			command.Env = append(command.Env, path)
		}

		return command.Run()
	}

	return errHelp
}

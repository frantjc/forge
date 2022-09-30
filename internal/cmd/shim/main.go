package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/frantjc/forge/pkg/envconv"
	"github.com/frantjc/forge/pkg/github/actions"
	"github.com/frantjc/go-js"
)

var (
	errHelp = errors.New("help")
	help    = fmt.Sprintf(`
%s [-c|-s|-e|-h] [args]

  -c   clone the given GitHub Action to the given path (default ".")
  -s   sleep
  -e   execute the given command after sourcing $GITHUB_PATH and $GITHUB_ENV
  -h   help

`, os.Args[0])
)

func main() {
	if err := mainE(); err != nil {
		if errors.Is(err, errHelp) {
			os.Stderr.WriteString(help)
			os.Exit(0)
		}

		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func mainE() error {
	var (
		ctx  = context.Background()
		args = os.Args
	)

	if len(args) < 2 {
		return errHelp
	}

	switch args[1] {
	// clone
	case "-c":
		if len(args) < 3 {
			return errHelp
		}

		var (
			usesStr = args[2]
			path    = "."
		)

		if len(args) > 3 {
			path = args[3]
		}

		parsed, err := actions.Parse(usesStr)
		if err != nil {
			return err
		}

		m, err := actions.Clone(ctx, parsed, &actions.CloneOpts{
			Path:     path,
			Insecure: true,
		})
		if err != nil {
			return err
		}

		return json.NewEncoder(os.Stdout).Encode(m)
	// sleep
	case "-s":
		os.Stdout.WriteString("zzz...")
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		close(sigs)
		_, err := os.Stdout.WriteString("\n")
		return err
	// exec
	case "-e":
		if len(args) < 3 {
			return errHelp
		}

		var (
			command        = exec.CommandContext(ctx, args[2], args[3:]...) //nolint:gosec
			githubEnvFile  = os.Getenv(actions.EnvVarEnv)
			githubPathFile = os.Getenv(actions.EnvVarPath)
		)

		command.Env = os.Environ()
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if githubEnv, err := envconv.ArrFromFile(githubEnvFile); err == nil {
			command.Env = append(command.Env, githubEnv...)
		} else {
			if _, err = os.Create(githubEnvFile); err != nil {
				return err
			}
		}

		var (
			path = "PATH=" + os.Getenv("PATH")
		)
		if runnerToolCache := os.Getenv(actions.EnvVarRunnerToolCache); runnerToolCache != "" {
			path += ":" + runnerToolCache
		}

		if githubPath, err := envconv.PathFromFile(githubPathFile); err == nil && githubPath != "" {
			path += ":" + githubPath
		} else {
			if _, err = os.Create(githubPathFile); err != nil {
				return err
			}
		}

		if i := js.FindIndex(command.Env, func(s string, _ int, _ []string) bool {
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

package forge

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/frantjc/forge/cloudbuild"
	"github.com/frantjc/forge/internal/bin"
	"github.com/frantjc/forge/internal/envconv"
	xos "github.com/frantjc/x/os"
	xslices "github.com/frantjc/x/slices"
)

type CloudBuild struct {
	cloudbuild.Step
}

func (o *CloudBuild) Run(ctx context.Context, containerRuntime ContainerRuntime, opts ...RunOpt) error {
	opt := runOptsWithDefaults(opts...)

	image, err := containerRuntime.PullImage(ctx, o.Name)
	if err != nil {
		return err
	}

	home := "/root"
	if config, err := image.Config(); err == nil {
		for _, env := range config.Env {
			if strings.HasPrefix(env, "HOME=") {
				home = strings.TrimPrefix(env, "HOME=")
				break
			}
		}
	}

	containerConfig, script, err := stepToContainerConfigAndScript(&o.Step, home, image, opt)
	if err != nil {
		return err
	}
	containerConfig.Mounts = overrideMounts(containerConfig.Mounts, opt.Mounts...)

	container, err := createSleepingContainer(ctx, containerRuntime, image, containerConfig, opt)
	if err != nil {
		return err
	}
	defer container.Stop(ctx) //nolint:errcheck
	// defer container.Remove(ctx) //nolint:errcheck

	if err = copyScriptToContainer(ctx, container, script, opt); err != nil {
		return err
	}

	if exitCode, err := container.Exec(ctx, containerConfig, opt.Streams); err != nil {
		return err
	} else if exitCode > 0 {
		return xos.NewExitCodeError(ErrContainerExitedWithNonzeroExitCode, exitCode)
	}

	return nil
}

func copyScriptToContainer(ctx context.Context, container Container, script string, opt *RunOpts) error {
	if !bin.HasShebang(script) {
		script = fmt.Sprintf("#!/usr/bin/env sh\nset -eo pipefail\n%s", script)
	}

	if err := container.CopyTo(ctx, opt.WorkingDir, bin.NewScriptTarArchive(script, ScriptName)); err != nil {
		return err
	}

	return nil
}

func stepToContainerConfigAndScript(step *cloudbuild.Step, home string, image Image, opt *RunOpts) (*ContainerConfig, string, error) {
	var (
		containerConfig = &ContainerConfig{
			Entrypoint: []string{},
			Env:        step.Env,
			WorkingDir: CloudBuildWorkingDir(opt.WorkingDir),
		}
		substitutions = step.Substitutions
		mapping       = func(s string) string {
			if substitution, ok := substitutions[s]; ok {
				return substitution
			}

			return fmt.Sprintf("$%s", s)
		}
		_home = os.Getenv("HOME")
	)
	if _home == "" {
		if u, err := user.Current(); err == nil {
			_home = u.HomeDir
		}
	}

	source := filepath.Join(_home, ".config/gcloud")
	if fi, err := os.Stat(source); err == nil && fi.IsDir() {
		containerConfig.Mounts = []Mount{
			{
				Source:      source,
				Destination: filepath.Join(home, ".config/gcloud"),
			},
		}
	}

	if step.Script != "" {
		if step.Entrypoint != "" || len(step.Args) > 0 {
			return nil, "", fmt.Errorf("cannot specify args or entrypoint with script")
		}

		containerConfig.Entrypoint = []string{filepath.Join(opt.WorkingDir, ScriptName)}
	} else {
		if lenArgs := len(step.Args); step.Entrypoint == "" || lenArgs == 0 {
			config, err := image.Config()
			if err != nil {
				return nil, "", err
			}

			if step.Entrypoint != "" {
				containerConfig.Entrypoint = []string{step.Entrypoint}
			} else {
				containerConfig.Entrypoint = config.Entrypoint
			}

			if lenArgs == 0 {
				containerConfig.Entrypoint = append(containerConfig.Entrypoint, config.Cmd...)
			} else {
				containerConfig.Entrypoint = append(containerConfig.Entrypoint, step.Args...)
			}
		} else {
			containerConfig.Entrypoint = append([]string{step.Entrypoint}, step.Args...)
		}
	}

	if step.DynamicSubstitutions {
		for range []byte{0, 0} {
			for k, v := range substitutions {
				substitutions[k] = os.Expand(v, mapping)
			}
		}
	}

	containerConfig.Entrypoint = xslices.Map(containerConfig.Entrypoint, func(s string, _ int) string {
		return os.Expand(s, mapping)
	})

	containerConfig.Env = xslices.Map(containerConfig.Env, func(s string, _ int) string {
		return os.Expand(s, mapping)
	})

	if step.AutomapSubstitutions {
		containerConfig.Env = append(containerConfig.Env, envconv.MapToArr(substitutions)...)
	}

	return containerConfig, os.Expand(step.Script, mapping), nil
}

package forgecloudbuild

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/cloudbuild"
	"github.com/frantjc/forge/envconv"
	"github.com/frantjc/forge/internal/bin"
	"github.com/frantjc/forge/internal/containerfs"
	xslice "github.com/frantjc/x/slice"
)

func StepToContainerConfigAndScript(step *cloudbuild.Step, home string, image forge.Image) (*forge.ContainerConfig, string, error) {
	return DefaultMapping.StepToContainerConfigAndScript(step, home, image)
}

func (m *Mapping) StepToContainerConfigAndScript(step *cloudbuild.Step, home string, image forge.Image) (*forge.ContainerConfig, string, error) {
	var (
		containerConfig = &forge.ContainerConfig{
			Entrypoint: []string{bin.ShimPath, "exec", "--sock", containerfs.ForgeSock, "--"},
			Env:        step.Env,
			WorkingDir: m.CloudBuildPath,
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
	if _, err := os.Stat(source); err == nil {
		containerConfig.Mounts = []forge.Mount{
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

		containerConfig.Cmd = bin.ScriptEntrypoint
	} else {
		if lenArgs := len(step.Args); step.Entrypoint == "" || lenArgs == 0 {
			config, err := image.Config()
			if err != nil {
				return nil, "", err
			}

			if step.Entrypoint != "" {
				containerConfig.Cmd = []string{step.Entrypoint}
			} else {
				containerConfig.Cmd = config.Entrypoint
			}

			if lenArgs == 0 {
				containerConfig.Cmd = append(containerConfig.Cmd, config.Cmd...)
			} else {
				containerConfig.Cmd = append(containerConfig.Cmd, step.Args...)
			}
		} else {
			containerConfig.Cmd = append([]string{step.Entrypoint}, step.Args...)
		}
	}

	if step.DynamicSubstitutions {
		for range []byte{0, 0} {
			for k, v := range substitutions {
				substitutions[k] = os.Expand(v, mapping)
			}
		}
	}

	containerConfig.Cmd = xslice.Map(containerConfig.Cmd, func(s string, _ int) string {
		return os.Expand(s, mapping)
	})

	containerConfig.Env = xslice.Map(containerConfig.Env, func(s string, _ int) string {
		return os.Expand(s, mapping)
	})

	if step.AutomapSubstitutions {
		containerConfig.Env = append(containerConfig.Env, envconv.MapToArr(substitutions)...)
	}

	return containerConfig, os.Expand(step.Script, mapping), nil
}

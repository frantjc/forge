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
	xslice "github.com/frantjc/x/slice"
)

func StepToContainerConfigAndScript(step *cloudbuild.Step, home string) (*forge.ContainerConfig, string, error) {
	return DefaultMapping.StepToContainerConfigAndScript(step, home)
}

func (m *Mapping) StepToContainerConfigAndScript(step *cloudbuild.Step, home string) (*forge.ContainerConfig, string, error) {
	var (
		containerConfig = &forge.ContainerConfig{
			Cmd:        step.Args,
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

	containerConfig.Mounts = []forge.Mount{
		{
			Source:      filepath.Join(_home, ".config/gcloud"),
			Destination: filepath.Join(home, ".config/gcloud"),
		},
	}

	if step.Script != "" {
		if step.Entrypoint != "" || len(step.Args) > 0 {
			return nil, "", fmt.Errorf("cannot specify args or entrypoint with script")
		}

		containerConfig.Entrypoint = bin.ScriptEntrypoint
	} else {
		containerConfig.Entrypoint = []string{step.Entrypoint}
	}

	if step.DynamicSubstitutions {
		for range []byte{0, 0} {
			for k, v := range substitutions {
				substitutions[k] = os.Expand(v, mapping)
			}
		}
	}

	containerConfig.Entrypoint = xslice.Map(containerConfig.Entrypoint, func(s string, _ int) string {
		return os.Expand(s, mapping)
	})

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

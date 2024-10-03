package forgeactions

import (
	"fmt"
	"path/filepath"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/envconv"
	"github.com/frantjc/forge/githubactions"
)

func ActionToConfigs(globalContext *githubactions.GlobalContext, uses *githubactions.Uses, with, environment map[string]string, actionMetadata *githubactions.Metadata, image forge.Image) ([]forge.ContainerConfig, error) {
	return DefaultMapping.ActionToConfigs(globalContext, uses, with, environment, actionMetadata, image)
}

func (m *Mapping) ActionToConfigs(globalContext *githubactions.GlobalContext, uses *githubactions.Uses, with, environment map[string]string, actionMetadata *githubactions.Metadata, image forge.Image) ([]forge.ContainerConfig, error) {
	containerConfigs := []forge.ContainerConfig{}
	globalContext = m.ConfigureGlobalContext(globalContext)

	if actionMetadata != nil {
		if actionMetadata.Runs != nil {
			actionDir, err := m.UsesToActionDirectory(uses)
			if err != nil {
				return nil, err
			}

			var (
				entrypoint = []string{}
				env        = append(envconv.MapToArr(environment), envconv.MapToArr(actionMetadata.Runs.Env)...)
				cmd        = actionMetadata.Runs.Args
				mounts     = []forge.Mount{
					{
						Source:      actionDir,
						Destination: m.ActionPath,
					},
					{
						Destination: m.Workspace,
					},
					{
						Destination: m.RunnerToolCache,
					},
					{
						Destination: m.RunnerTemp,
					},
				}
				entrypoints [][]string
			)

			switch actionMetadata.Runs.Using {
			case githubactions.RunsUsingNode12, githubactions.RunsUsingNode16, githubactions.RunsUsingNode20:
				entrypoint = []string{"node"}

				if pre := actionMetadata.Runs.Pre; pre != "" {
					entrypoints = append(entrypoints, []string{filepath.Join(m.ActionPath, pre)})
				}

				if main := actionMetadata.Runs.Main; main != "" {
					entrypoints = append(entrypoints, []string{filepath.Join(m.ActionPath, main)})
				}
			case githubactions.RunsUsingDocker:
				if pre := actionMetadata.Runs.PreEntrypoint; pre != "" {
					entrypoints = append(entrypoints, []string{pre})
				}

				if main := actionMetadata.Runs.Entrypoint; main != "" {
					entrypoints = append(entrypoints, []string{main})
				} else {
					config, err := image.Config()
					if err != nil {
						return nil, err
					}

					entrypoints = append(entrypoints, config.Entrypoint)
				}
			default:
				return nil, fmt.Errorf("unsupported runs using: %s", actionMetadata.Runs.Using)
			}

			unexpandedInputs, err := actionMetadata.InputsFromWith(with)
			if err != nil {
				return nil, err
			}

			var (
				inputs   = make(map[string]string, len(unexpandedInputs))
				expander = githubactions.ExpandFunc(globalContext.GetString)
			)
			for k, v := range unexpandedInputs {
				inputs[k] = expander.ExpandString(v)
			}

			globalContext.InputsContext = inputs
			env = append(env, globalContext.Env()...)
			env = append(env,
				fmt.Sprintf("%s=%s", githubactions.EnvVarPath, m.GitHubPathPath),
				fmt.Sprintf("%s=%s", githubactions.EnvVarEnv, m.GitHubEnvPath),
				fmt.Sprintf("%s=%s", githubactions.EnvVarOutput, m.GitHubOutputPath),
				fmt.Sprintf("%s=%s", githubactions.EnvVarState, m.GitHubStatePath),
			)

			for _, ep := range entrypoints {
				if len(ep) > 0 {
					if len(entrypoint) > 0 {
						containerConfigs = append(containerConfigs, forge.ContainerConfig{
							Entrypoint: entrypoint,
							Cmd:        append(ep, cmd...),
							Env:        env,
							Mounts:     mounts,
							WorkingDir: m.Workspace,
						})
					} else {
						containerConfigs = append(containerConfigs, forge.ContainerConfig{
							Entrypoint: ep,
							Cmd:        cmd,
							Env:        env,
							Mounts:     mounts,
							WorkingDir: m.Workspace,
						})
					}
				}
			}
		}
	}

	return containerConfigs, nil
}

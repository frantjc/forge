package forgeactions

import (
	"path/filepath"
	"strings"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/bin"
	"github.com/frantjc/forge/pkg/envconv"
	"github.com/frantjc/forge/pkg/githubactions"
)

func ActionToConfigs(globalContext *githubactions.GlobalContext, uses *githubactions.Uses, with, environment map[string]string, actionMetadata *githubactions.Metadata, image forge.Image) ([]*forge.ContainerConfig, error) {
	return DefaultMapping.ActionToConfigs(globalContext, uses, with, environment, actionMetadata, image)
}

func (m *Mapping) ActionToConfigs(globalContext *githubactions.GlobalContext, uses *githubactions.Uses, with, environment map[string]string, actionMetadata *githubactions.Metadata, image forge.Image) ([]*forge.ContainerConfig, error) {
	var (
		_                = forge.NewLogger()
		containerConfigs = []*forge.ContainerConfig{}
	)
	globalContext = m.ConfigureGlobalContext(globalContext)

	if actionMetadata != nil {
		if actionMetadata.GetRuns() != nil {
			actionDir, err := m.UsesToActionDirectory(uses)
			if err != nil {
				return nil, err
			}

			var (
				entrypoint = []string{bin.ShimPath, "-e"}
				env        = append(envconv.MapToArr(environment), envconv.MapToArr(actionMetadata.GetRuns().GetEnv())...)
				cmd        = actionMetadata.GetRuns().GetArgs()
				mounts     = []*forge.Mount{
					{
						Source:      actionDir,
						Destination: m.GetActionPath(),
					},
					{
						Destination: m.GetWorkspace(),
					},
					{
						Destination: m.GetRunnerToolCache(),
					},
					{
						Destination: m.GetRunnerTemp(),
					},
					{
						Destination: m.GetGitHubPath(),
					},
				}
				entrypoints []string
			)

			switch actionMetadata.GetRuns().GetUsing() {
			case githubactions.RunsUsingNode12, githubactions.RunsUsingNode16:
				entrypoint = append(entrypoint, "node")
				if pre := actionMetadata.GetRuns().GetPre(); pre != "" {
					entrypoints = append(entrypoints, filepath.Join(m.GetActionPath(), pre))
				}

				if main := actionMetadata.GetRuns().GetMain(); main != "" {
					entrypoints = append(entrypoints, filepath.Join(m.GetActionPath(), main))
				}
			case githubactions.RunsUsingDocker:
				if pre := actionMetadata.GetRuns().GetPreEntrypoint(); pre != "" {
					entrypoints = append(entrypoints, pre)
				}

				if main := actionMetadata.GetRuns().GetMain(); main != "" {
					entrypoints = append(entrypoints, main)
				} else {
					config, err := image.Config()
					if err != nil {
						return nil, err
					}

					entrypoints = append(entrypoints, strings.Join(config.Entrypoint, " "))
				}
			}

			unexpandedInputs, err := actionMetadata.InputsFromWith(with)
			if err != nil {
				return nil, err
			}

			var (
				inputs   = make(map[string]string, len(unexpandedInputs))
				expander = githubactions.NewExpander(globalContext.GetString)
			)
			for k, v := range unexpandedInputs {
				inputs[k] = expander.ExpandString(v)
			}

			globalContext.InputsContext = inputs
			env = append(env, globalContext.Env()...)
			env = append(env, githubactions.EnvVarPath+"="+m.GetGitHubPathPath(), githubactions.EnvVarEnv+"="+m.GetGitHubEnvPath())

			for _, s := range entrypoints {
				if s != "" {
					containerConfigs = append(containerConfigs, &forge.ContainerConfig{
						Entrypoint: entrypoint,
						Cmd:        append([]string{s}, cmd...),
						Env:        env,
						Mounts:     mounts,
					})
				}
			}
		}
	}

	return containerConfigs, nil
}

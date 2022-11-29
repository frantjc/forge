package forgeactions

import "github.com/frantjc/forge/pkg/githubactions"

func ConfigureGlobalContext(globalContext *githubactions.GlobalContext) *githubactions.GlobalContext {
	return DefaultMapping.ConfigureGlobalContext(globalContext)
}

func (m *Mapping) ConfigureGlobalContext(globalContext *githubactions.GlobalContext) *githubactions.GlobalContext {
	if globalContext == nil {
		globalContext = githubactions.NewGlobalContextFromEnv()
	}

	if globalContext.GitHubContext == nil {
		globalContext.GitHubContext = &githubactions.GitHubContext{}
	}

	globalContext.GitHubContext.Workspace = m.Workspace
	globalContext.GitHubContext.ActionPath = m.ActionPath

	if globalContext.RunnerContext == nil {
		globalContext.RunnerContext = &githubactions.RunnerContext{}
	}

	globalContext.RunnerContext.Temp = m.RunnerTemp
	globalContext.RunnerContext.ToolCache = m.RunnerToolCache

	return globalContext
}

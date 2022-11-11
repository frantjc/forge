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

	globalContext.GitHubContext.Workspace = m.GetWorkspace()
	globalContext.GitHubContext.ActionPath = m.GetActionPath()

	if globalContext.RunnerContext == nil {
		globalContext.RunnerContext = &githubactions.RunnerContext{}
	}

	globalContext.RunnerContext.Temp = m.GetRunnerTemp()
	globalContext.RunnerContext.ToolCache = m.GetRunnerToolCache()

	return globalContext
}

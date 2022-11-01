package forgeactions

import "github.com/frantjc/forge/pkg/github/actions"

func ConfigureGlobalContext(globalContext *actions.GlobalContext) *actions.GlobalContext {
	return DefaultMapping.ConfigureGlobalContext(globalContext)
}

func (m *Mapping) ConfigureGlobalContext(globalContext *actions.GlobalContext) *actions.GlobalContext {
	if globalContext == nil {
		globalContext = actions.NewGlobalContextFromEnv()
	}

	if globalContext.GitHubContext == nil {
		globalContext.GitHubContext = &actions.GitHubContext{}
	}

	globalContext.GitHubContext.Workspace = m.Workspace
	globalContext.GitHubContext.ActionPath = m.ActionPath

	if globalContext.RunnerContext == nil {
		globalContext.RunnerContext = &actions.RunnerContext{}
	}

	globalContext.RunnerContext.Temp = m.RunnerTemp
	globalContext.RunnerContext.ToolCache = m.RunnerToolCache

	return globalContext
}

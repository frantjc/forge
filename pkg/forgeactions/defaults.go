package fa

import cfs "github.com/frantjc/forge/internal/containerfs"

var (
	DefaultMapping = &Mapping{
		Workspace:       DefaultWorkspace,
		ActionPath:      DefaultActionPath,
		RunnerTemp:      DefaultRunnerTemp,
		RunnerToolCache: DefaultRunnerToolCache,
		GitHubPath:      DefaultGitHubPath,
		GitHubPathPath:  DefaultGitHubPathPath,
		GitHubEnvPath:   DefaultGitHubEnvPath,
	}
)

var (
	DefaultRootPath        = cfs.WorkingDir
	DefaultWorkspace       = DefaultRootPath + "/workspace"
	DefaultActionPath      = DefaultRootPath + "/action"
	DefaultRunnerPath      = DefaultRootPath + "/runner"
	DefaultRunnerTemp      = DefaultRunnerPath + "/tmp"
	DefaultRunnerToolCache = DefaultRunnerPath + "/toolcache"
	DefaultGitHubPath      = DefaultRootPath + "/github"
	DefaultGitHubPathPath  = DefaultGitHubPath + "/path.txt"
	DefaultGitHubEnvPath   = DefaultGitHubPath + "/env.txt"
)

package actions2container

var (
	DefaultMap = &Map{
		Workspace:                 DefaultWorkspace,
		ActionPath:                DefaultActionPath,
		RunnerTemp:                DefaultRunnerTemp,
		RunnerToolCache:           DefaultRunnerToolCache,
		GitHubPath:                DefaultGitHubPath,
		GitHubPathPath:            DefaultGitHubPathPath,
		GitHubEnvPath:             DefaultGitHubEnvPath,
		RunnerToolCacheVolumeName: DefaultRunnerToolCacheVolumeName,
	}
)

const (
	DefaultRootPath        = "/out"
	DefaultWorkspace       = DefaultRootPath + "/workspace"
	DefaultActionPath      = DefaultRootPath + "/action"
	DefaultRunnerTemp      = DefaultRootPath + "/runner/tmp"
	DefaultRunnerToolCache = DefaultRootPath + "/runner/toolcache"
	DefaultGitHubPath      = DefaultRootPath + "/github"
	DefaultGitHubPathPath  = DefaultGitHubPath + "/path.txt"
	DefaultGitHubEnvPath   = DefaultGitHubPath + "/env.txt"
)

const (
	DefaultRunnerToolCacheVolumeName = "runner-cache"
)

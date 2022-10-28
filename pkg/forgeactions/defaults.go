package fa

import (
	"github.com/adrg/xdg"
	"github.com/frantjc/forge"
)

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
		ActionCache:               DefaultActionCache,
	}
)

const (
	DefaultRootPath        = forge.WorkingDir
	DefaultWorkspace       = DefaultRootPath + "/workspace"
	DefaultActionPath      = DefaultRootPath + "/action"
	DefaultRunnerPath      = DefaultRootPath + "/runner"
	DefaultRunnerTemp      = DefaultRunnerPath + "/tmp"
	DefaultRunnerToolCache = DefaultRunnerPath + "/toolcache"
	DefaultGitHubPath      = DefaultRootPath + "/github"
	DefaultGitHubPathPath  = DefaultGitHubPath + "/path.txt"
	DefaultGitHubEnvPath   = DefaultGitHubPath + "/env.txt"
)

const (
	DefaultRunnerToolCacheVolumeName = "runner-cache"
)

var (
	DefaultActionCache = xdg.CacheHome
)

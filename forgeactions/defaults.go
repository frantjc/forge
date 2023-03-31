package forgeactions

import "github.com/frantjc/forge/internal/containerfs"

var DefaultMapping = &Mapping{
	Workspace:        DefaultWorkspace,
	ActionPath:       DefaultActionPath,
	RunnerTemp:       DefaultRunnerTemp,
	RunnerToolCache:  DefaultRunnerToolCache,
	GitHubPath:       DefaultGitHubPath,
	GitHubPathPath:   DefaultGitHubPathPath,
	GitHubEnvPath:    DefaultGitHubEnvPath,
	GitHubOutputPath: DefaultGitHubOutputPath,
	GitHubStatePath:  DefaultGitHubStatePath,
}

var (
	DefaultRootPath         = containerfs.WorkingDir
	DefaultWorkspace        = DefaultRootPath + "/workspace"
	DefaultActionPath       = DefaultRootPath + "/action"
	DefaultRunnerPath       = DefaultRootPath + "/runner"
	DefaultRunnerTemp       = DefaultRunnerPath + "/tmp"
	DefaultRunnerToolCache  = DefaultRunnerPath + "/toolcache"
	DefaultGitHubPath       = DefaultRootPath + "/github"
	DefaultGitHubPathPath   = DefaultGitHubPath + "/add_path"
	DefaultGitHubEnvPath    = DefaultGitHubPath + "/set_env"
	DefaultGitHubOutputPath = DefaultGitHubPath + "/set_output"
	DefaultGitHubStatePath  = DefaultGitHubPath + "/save_state"
)

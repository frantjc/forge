package forge

import (
	"path/filepath"
)

func CloudBuildWorkingDir(workingDir string) string {
	return filepath.Join(workingDir, "cloudbuild")
}

func InterceptingDockerSock(workingDir string) string {
	return filepath.Join(workingDir, "forge.sock")
}

func GitHubWorkspace(workingDir string) string {
	return filepath.Join(workingDir, "github/workspace")
}

func GitHubActionPath(workingDir string) string {
	return filepath.Join(workingDir, "github/action")
}

func GitHubRunnerTmp(workingDir string) string {
	return filepath.Join(workingDir, "github/runner/tmp")
}

func GitHubRunnerToolCache(workingDir string) string {
	return filepath.Join(workingDir, "github/runner/toolcache")
}

func GitHubPath(workingDir string) string {
	return filepath.Join(workingDir, "github/add_path")
}

func GitHubEnv(workingDir string) string {
	return filepath.Join(workingDir, "github/set_env")
}

func GitHubOutput(workingDir string) string {
	return filepath.Join(workingDir, "github/set_output")
}

func GitHubState(workingDir string) string {
	return filepath.Join(workingDir, "github/save_state")
}

func ConcourseResourceWorkingDir(workingDir string) string {
	return filepath.Join(workingDir, "concourse/resource")
}

func AzureDevOpsTaskWorkingDir(workingDir string) string {
	return filepath.Join(workingDir, "task")
}

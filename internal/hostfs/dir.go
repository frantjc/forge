package hostfs

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

// CacheHome is the directory on the host machine
// where Forge caches stuff.
var CacheHome = filepath.Join(xdg.CacheHome, "forge")

var (
	gitHubActionsHome = filepath.Join(CacheHome, "/github/actions")
	// ActionsCache is the directory on the host machine where
	// all GitHub Action repositories are stored.
	ActionsCache = filepath.Join(gitHubActionsHome, "/actions")
	// RunnerTmp is the the directory on the host machine used as the source
	// for the mount at RUNNER_TEMP.
	RunnerTmp = filepath.Join(gitHubActionsHome, "/runner/tmp")
	// RunnerTmp is the the directory on the host machine used as the source
	// for the mount at RUNNER_TOOLCACHE.
	RunnerToolCache = filepath.Join(gitHubActionsHome, "/runner/toolcache")
)

var CloudBuildWorkspace = filepath.Join(CacheHome, "/cloudbuild/workspace")

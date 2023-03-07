package hostfs

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

// CacheHome is the directory on the host machine
// where Forge caches stuff.
var CacheHome = filepath.Join(xdg.CacheHome, ".forge")

var (
	// ActionsCache is the directory on the host machine where
	// all GitHub Action repositories are stored.
	ActionsCache = filepath.Join(CacheHome, "/actions")
	// RunnerTmp is the the directory on the host machine used as the source
	// for the mount at RUNNER_TEMP.
	RunnerTmp = filepath.Join(ActionsCache, "/runner/tmp")
	// RunnerTmp is the the directory on the host machine used as the source
	// for the mount at RUNNER_TOOLCACHE.
	RunnerToolCache = filepath.Join(ActionsCache, "/runner/toolcache")
)

var (
	CircleCICache = filepath.Join(CacheHome, "/circleci")
	CircleCIHome  = filepath.Join(CircleCICache, "/home")
)

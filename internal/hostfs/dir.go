package hostfs

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

var (
	// ActionsCache is the directory on the host machine where
	// all GitHub Actions-related stuff is stored.
	ActionsCache = filepath.Join(xdg.CacheHome, ".forge")
	// ActionCache is the directory on the host machine where
	// all GitHub Action repositories are stored.
	ActionCache = filepath.Join(ActionsCache, "/actions")
	// RunnerTmp is the the directory on the host machine used as the source
	// for the mount at RUNNER_TEMP.
	RunnerTmp = filepath.Join(ActionsCache, "/runner/tmp")
	// RunnerTmp is the the directory on the host machine used as the source
	// for the mount at RUNNER_TOOLCACHE.
	RunnerToolcache = filepath.Join(ActionsCache, "/runner/toolcache")
)

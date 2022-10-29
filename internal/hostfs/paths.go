package hfs

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

var (
	ActionCache = filepath.Join(xdg.CacheHome, ".forge")
	// RunnerTmp is the the directory on the host machine used as the source
	// for the mount at RUNNER_TEMP.
	RunnerTmp = filepath.Join(ActionCache, "/runner/tmp")
	// RunnerTmp is the the directory on the host machine used as the source
	// for the mount at RUNNER_TOOLCACHE.
	RunnerToolcache = filepath.Join(ActionCache, "/runner/toolcache")
)

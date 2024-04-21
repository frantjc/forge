package bin

import (
	"path/filepath"

	"github.com/frantjc/forge/internal/containerfs"
)

var (
	ShimPath       = filepath.Join(containerfs.WorkingDir, ShimName)
	ShimEntrypoint = []string{ShimPath}
)

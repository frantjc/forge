package containerfs

import (
	"path/filepath"

	"github.com/google/uuid"
)

var (
	// WorkingDir is the directory Ores are ran from the context of.
	WorkingDir = filepath.Join("/forge", uuid.NewString())
)

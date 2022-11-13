package containerfs

import (
	"path/filepath"

	"github.com/google/uuid"
)

// WorkingDir is the directory Ores are ran from the context of.
var WorkingDir = filepath.Join("/forge", uuid.NewString())

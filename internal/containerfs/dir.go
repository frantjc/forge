package containerfs

import (
	"github.com/google/uuid"
)

var (
	// WorkingDir is the directory Ores are ran from the context of.
	WorkingDir = "/forge/" + uuid.NewString()
)

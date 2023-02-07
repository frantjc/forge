package containerfs

import "github.com/google/uuid"

// WorkingDir is the directory in the containers forge creates
// where extra stuff such as forge's shim, docker.sock and
// GitHub Action repositories gets mounted to.
var WorkingDir = "/" + uuid.NewString()

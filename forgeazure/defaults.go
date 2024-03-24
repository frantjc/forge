package forgeazure

import "github.com/frantjc/forge/internal/containerfs"

var (
	DefaultTaskPath = containerfs.WorkingDir + "/task"
	DefaultMapping  = &Mapping{
		TaskPath: DefaultTaskPath,
	}
)

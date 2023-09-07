package forgeazure

import "github.com/frantjc/forge/internal/containerfs"

var DefaultMapping = &Mapping{
	TaskPath: DefaultTaskPath,
}

var (
	DefaultRootPath = containerfs.WorkingDir
	DefaultTaskPath = DefaultRootPath + "/task"
)

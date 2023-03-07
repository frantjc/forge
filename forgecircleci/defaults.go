package forgecircleci

import "github.com/frantjc/forge/internal/containerfs"

var DefaultMapping = &Mapping{
	Home: DefaultHome,
}

var (
	DefaultRootPath = containerfs.WorkingDir
	DefaultHome     = DefaultRootPath + "/home"
)

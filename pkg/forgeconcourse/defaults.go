package forgeconcourse

import "github.com/frantjc/forge/internal/containerfs"

var (
	DefaultMapping = &Mapping{
		RootPath: DefaultRootPath,
	}
)

var (
	DefaultRootPath = containerfs.WorkingDir
)

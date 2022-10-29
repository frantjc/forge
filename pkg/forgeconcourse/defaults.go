package fc

import cfs "github.com/frantjc/forge/internal/containerfs"

var (
	DefaultMapping = &Mapping{
		RootPath: DefaultRootPath,
	}
)

var (
	DefaultRootPath = cfs.WorkingDir
)

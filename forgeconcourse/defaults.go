package forgeconcourse

import "github.com/frantjc/forge/internal/containerfs"

var (
	DefaultRootPath = containerfs.WorkingDir + "/resource"
	DefaultMapping  = &Mapping{
		RootPath: DefaultRootPath,
	}
)

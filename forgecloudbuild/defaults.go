package forgecloudbuild

import "github.com/frantjc/forge/internal/containerfs"

var (
	DefaultCloudBuildPath = containerfs.WorkingDir + "/cloudbuild"
	DefaultMapping        = &Mapping{
		CloudBuildPath: DefaultCloudBuildPath,
	}
)

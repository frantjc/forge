package forge

import "errors"

var (
	// ErrContainerExecutedWithNonzeroExitCode is returned when a container ran by
	// an Ore returned an unexpected nonzero exit code. Used with errors.Is to avoid
	// printing this error's text so that os.Exit doesn't double-inform a user of
	// the same nonzero exit code.
	ErrContainerExitedWithNonzeroExitCode = errors.New("container exited with nonzero exit code")
	// ErrCannotBuildDockerfile will be returned when a forge.ContainerRuntime
	// does not implement ImageBuilder.
	ErrCannotBuildDockerfile = errors.New("runtime cannot build Dockerfile")
)

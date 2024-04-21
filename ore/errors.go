package ore

import "errors"

var ErrContainerExitedWithNonzeroExitCode = errors.New("container exited with nonzero exit code")

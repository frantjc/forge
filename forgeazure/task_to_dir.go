package forgeazure

import (
	"path/filepath"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/azuredevops"
)

// TaskReferenceToDirectory is a re-export of DefaultMapping.TaskReferenceToDirectory
// for convenience purposes.
func TaskReferenceToDirectory(ref *azuredevops.TaskReference) (string, error) {
	return DefaultMapping.TaskReferenceToDirectory(ref)
}

// TaskReferenceToDirectory takes an *azuredevops.TaskReference and returns the path to the
// directory where the corresponding task.json can be found on the host machine.
func (m *Mapping) TaskReferenceToDirectory(ref *azuredevops.TaskReference) (string, error) {
	if ref.IsLocal() {
		return filepath.Abs(ref.Path)
	}

	return "", forge.ErrUnimplemented
}

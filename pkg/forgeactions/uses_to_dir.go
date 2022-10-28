package fa

import (
	"path/filepath"

	"github.com/frantjc/forge/pkg/github/actions"
)

func UsesToDirectory(uses *actions.Uses) (string, error) {
	return DefaultMap.UsesToDirectory(uses)
}

func (m *Map) UsesToDirectory(uses *actions.Uses) (string, error) {
	if uses.IsLocal() {
		return filepath.Abs(uses.GetPath())
	}

	return filepath.Join(m.GetActionCache(), uses.GetPath(), uses.GetVersion()), nil
}

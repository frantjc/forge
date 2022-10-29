package fa

import (
	"path/filepath"

	hfs "github.com/frantjc/forge/internal/hostfs"
	"github.com/frantjc/forge/pkg/github/actions"
)

func UsesToRootDirectory(uses *actions.Uses) (string, error) {
	return DefaultMapping.UsesToRootDirectory(uses)
}

func (m *Mapping) UsesToRootDirectory(uses *actions.Uses) (string, error) {
	if uses.IsLocal() {
		return filepath.Abs(uses.GetPath())
	}

	return filepath.Join(hfs.ActionCache, uses.GetRepository(), uses.GetVersion()), nil
}

func UsesToActionDirectory(uses *actions.Uses) (string, error) {
	return DefaultMapping.UsesToActionDirectory(uses)
}

func (m *Mapping) UsesToActionDirectory(uses *actions.Uses) (string, error) {
	if uses.IsLocal() {
		return m.UsesToRootDirectory(uses)
	}

	return filepath.Join(hfs.ActionCache, uses.GetRepository(), uses.GetVersion(), uses.GetActionPath()), nil
}

package forgeactions

import (
	"path/filepath"

	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/internal/hostfs"
)

func UsesToRootDirectory(uses *githubactions.Uses) (string, error) {
	return DefaultMapping.UsesToRootDirectory(uses)
}

func (m *Mapping) UsesToRootDirectory(uses *githubactions.Uses) (string, error) {
	if uses.IsLocal() {
		return filepath.Abs(uses.Path)
	}

	return filepath.Join(hostfs.ActionsCache, uses.GetRepository(), uses.Version), nil
}

func UsesToActionDirectory(uses *githubactions.Uses) (string, error) {
	return DefaultMapping.UsesToActionDirectory(uses)
}

func (m *Mapping) UsesToActionDirectory(uses *githubactions.Uses) (string, error) {
	if uses.IsLocal() {
		return m.UsesToRootDirectory(uses)
	}

	return filepath.Join(hostfs.ActionsCache, uses.GetRepository(), uses.Version, uses.GetActionPath()), nil
}

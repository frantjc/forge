package forgeactions

import (
	"path/filepath"

	"github.com/frantjc/forge/githubactions"
	"github.com/frantjc/forge/internal/hostfs"
)

// UsesToRootDirectory is a re-export of DefaultMapping.UsesToRootDirectory
// for convenience purposes.
func UsesToRootDirectory(uses *githubactions.Uses) (string, error) {
	return DefaultMapping.UsesToRootDirectory(uses)
}

// UsesToRootDirectory takes a *githubactions.Uses and returns the path to
// where the corresponding git repository can be found on the host machine.
func (m *Mapping) UsesToRootDirectory(uses *githubactions.Uses) (string, error) {
	if uses.IsLocal() {
		return filepath.Abs(uses.Path)
	}

	return filepath.Join(hostfs.ActionsCache, uses.GetRepository(), uses.Version), nil
}

// UsesToActionDirectory is a re-export of DefaultMapping.UsesToActionDirectory
// for convenience purposes.
func UsesToActionDirectory(uses *githubactions.Uses) (string, error) {
	return DefaultMapping.UsesToActionDirectory(uses)
}

// UsesToRootDirectory takes a *githubactions.Uses and returns the path to the
// directory where the corresponding action.yml can be found on the host machine.
func (m *Mapping) UsesToActionDirectory(uses *githubactions.Uses) (string, error) {
	if uses.IsLocal() {
		return m.UsesToRootDirectory(uses)
	}

	return filepath.Join(hostfs.ActionsCache, uses.GetRepository(), uses.Version, uses.GetActionPath()), nil
}

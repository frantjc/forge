package forgeactions

import (
	"context"

	"github.com/frantjc/forge/githubactions"
)

// GetUsesMetadata is a re-export of DefaultMapping.GetUsesMetadata
// for convenience purposes.
func GetUsesMetadata(ctx context.Context, uses *githubactions.Uses) (*githubactions.Metadata, error) {
	return DefaultMapping.GetUsesMetadata(ctx, uses)
}

// GetUsesMetadata gets the action.yml for the given *githubactions.Uses.
func (m *Mapping) GetUsesMetadata(ctx context.Context, uses *githubactions.Uses) (*githubactions.Metadata, error) {
	dir, err := m.UsesToRootDirectory(uses)
	if err != nil {
		return nil, err
	}

	return githubactions.GetUsesMetadata(ctx, uses, dir)
}

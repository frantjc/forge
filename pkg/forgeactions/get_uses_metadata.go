package forgeactions

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/githubactions"
)

func GetUsesMetadata(ctx context.Context, uses *githubactions.Uses) (*githubactions.Metadata, error) {
	return DefaultMapping.GetUsesMetadata(ctx, uses)
}

func (m *Mapping) GetUsesMetadata(ctx context.Context, uses *githubactions.Uses) (*githubactions.Metadata, error) {
	_ = forge.LoggerFrom(ctx)

	dir, err := m.UsesToRootDirectory(uses)
	if err != nil {
		return nil, err
	}

	return githubactions.GetUsesMetadata(ctx, uses, dir)
}

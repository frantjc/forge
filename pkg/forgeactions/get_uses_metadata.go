package fa

import (
	"context"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/github/actions"
)

func GetUsesMetadata(ctx context.Context, uses *actions.Uses) (*actions.Metadata, error) {
	return DefaultMapping.GetUsesMetadata(ctx, uses)
}

func (m *Mapping) GetUsesMetadata(ctx context.Context, uses *actions.Uses) (*actions.Metadata, error) {
	_ = forge.LoggerFrom(ctx)

	dir, err := m.UsesToRootDirectory(uses)
	if err != nil {
		return nil, err
	}

	return actions.GetUsesMetadata(ctx, uses, dir)
}

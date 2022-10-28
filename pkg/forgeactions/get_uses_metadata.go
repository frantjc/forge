package fa

import (
	"context"

	"github.com/frantjc/forge/pkg/github/actions"
)

func GetUsesMetadata(ctx context.Context, uses *actions.Uses) (*actions.Metadata, error) {
	return DefaultMap.GetUsesMetadata(ctx, uses)
}

func (m *Map) GetUsesMetadata(ctx context.Context, uses *actions.Uses) (*actions.Metadata, error) {
	dir, err := m.UsesToDirectory(uses)
	if err != nil {
		return nil, err
	}

	return actions.GetUsesMetadata(ctx, uses, dir)
}

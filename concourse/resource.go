package concourse

import "context"

type resourcesKey struct{}

func WithResources(ctx context.Context, resources ...*Resource) context.Context {
	return context.WithValue(ctx, resourcesKey{}, resources)
}

func ResourcesFrom(ctx context.Context) (resources []*Resource, ok bool) {
	resources, ok = ctx.Value(resourcesKey{}).([]*Resource)
	return
}

type Resource struct {
	Name   string            `json:"name,omitempty"`
	Type   string            `json:"type,omitempty"`
	Source map[string]string `json:"source,omitempty"`
}

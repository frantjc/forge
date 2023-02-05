package concourse

import "context"

type resourcesKey struct{}

// WithResources embeds the given resources into the given context.
func WithResources(ctx context.Context, resources ...*Resource) context.Context {
	return context.WithValue(ctx, resourcesKey{}, resources)
}

// ResourcesFrom returns whatever resources have been embedded in
// given context.
func ResourcesFrom(ctx context.Context) (resources []*Resource, ok bool) {
	resources, ok = ctx.Value(resourcesKey{}).([]*Resource)
	return
}

// Resource is the struct which has the YAML encoding of a resource
// as it would appear in the `resource` array of a Concourse pipeline
// configuration file.
type Resource struct {
	Name   string            `json:"name,omitempty"`
	Type   string            `json:"type,omitempty"`
	Source map[string]string `json:"source,omitempty"`
}

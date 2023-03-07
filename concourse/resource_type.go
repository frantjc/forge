package concourse

import "context"

type resourceTypesKey struct{}

// WithResourceTypes embeds the given resource types into the given context.
func WithResourceTypes(ctx context.Context, resourceTypes ...*ResourceType) context.Context {
	return context.WithValue(ctx, resourceTypesKey{}, resourceTypes)
}

// ResourceTypesFrom returns whatever resource types have been embedded in
// given context.
func ResourceTypesFrom(ctx context.Context) (resourceTypes []ResourceType, ok bool) {
	resourceTypes, ok = ctx.Value(resourceTypesKey{}).([]ResourceType)
	return
}

// ResourceType is the struct which has the YAML encoding of a resource type
// as it would appear in the `resource_types` array of a Concourse pipeline
// configuration file.
type ResourceType struct {
	Name       string `json:"name,omitempty"`
	Privileged bool   `json:"privileged,omitempty"`
	Source     *struct {
		Repository string `json:"repository,omitempty"`
		Tag        string `json:"tag,omitempty"`
	} `json:"source,omitempty"`
}

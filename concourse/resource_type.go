package concourse

import "context"

type resourceTypesKey struct{}

func WithResourceTypes(ctx context.Context, resourceTypes ...*ResourceType) context.Context {
	return context.WithValue(ctx, resourceTypesKey{}, resourceTypes)
}

func ResourceTypesFrom(ctx context.Context) (resourceTypes []*ResourceType, ok bool) {
	resourceTypes, ok = ctx.Value(resourceTypesKey{}).([]*ResourceType)
	return
}

type ResourceType struct {
	Name       string  `json:"name,omitempty"`
	Source     *Source `json:"source,omitempty"`
	Privileged bool    `json:"privileged,omitempty"`
}

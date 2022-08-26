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

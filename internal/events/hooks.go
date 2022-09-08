package events

import "context"

type Hookable interface {
	Invoke(*Event)
}

type Hook func(context.Context, *Event)

func (h Hook) Invoke(ctx context.Context, event *Event) {
	h(ctx, event)
}

type Hooks []Hook

func (s Hooks) Invoke(ctx context.Context, event *Event) {
	for _, h := range s {
		h.Invoke(ctx, event)
	}
}

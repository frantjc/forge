package events

import "context"

type EventStream interface {
	Emit(context.Context, *Event)
	Listen(context.Context, Type, ...Hook)
}

type Stream map[Type]Hooks

func (s Stream) Emit(ctx context.Context, event *Event) {
	for eventType, hooks := range s {
		if eventType == Type(event.Type) {
			hooks.Invoke(ctx, event)
		}
	}
}

func (s Stream) Listen(ctx context.Context, eventType Type, hooks ...Hook) {
	if streamHooks, ok := s[eventType]; ok {
		s[eventType] = append(streamHooks, hooks...)
	} else {
		s[eventType] = hooks
	}
}

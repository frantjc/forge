package events

import "context"

var (
	DefaultStream = Stream{}
)

func Emit(ctx context.Context, event *Event) {
	DefaultStream.Emit(ctx, event)
}

func Listen(ctx context.Context, eventType Type, hooks ...Hook) {
	DefaultStream.Listen(ctx, eventType, hooks...)
}

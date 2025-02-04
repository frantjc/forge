package forge

import (
	"context"
	"sync"
)

var HookContainerStarted = new(Hook[Container])

type Hook[T any] struct {
	Listeners []func(context.Context, T)
	sync.Mutex
}

func (h *Hook[T]) Dispatch(ctx context.Context, t T) {
	h.Lock()
	defer h.Unlock()

	for _, l := range h.Listeners {
		l(ctx, t)
	}
}

func (h *Hook[T]) Listen(f ...func(context.Context, T)) {
	h.Lock()
	defer h.Unlock()

	h.Listeners = append(h.Listeners, f...)
}

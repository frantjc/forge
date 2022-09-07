package events

type Hookable interface {
	Invoke(*Event)
}

type Hook func(*Event)

func (h Hook) Invoke(event *Event) {
	h(event)
}

type Hooks []Hook

func (s Hooks) Invoke(event *Event) {
	for _, h := range s {
		h.Invoke(event)
	}
}

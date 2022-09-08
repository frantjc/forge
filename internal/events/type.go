package events

const (
	ContainerCreated Type = "container.created"
)

type Type string

func (t Type) String() string {
	return string(t)
}

func (t Type) GoString() string {
	return "Type(" + t.String() + ")"
}

package events

func NewContainerMetadata(id string) ContainerMetadata {
	m := ContainerMetadata{}
	m["id"] = id
	return m
}

type ContainerMetadata map[string]string

func (m ContainerMetadata) GetId() string { //nolint:revive // matching .proto style
	return m["id"]
}

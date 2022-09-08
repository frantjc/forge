package events

import "github.com/frantjc/forge"

func NewContainerMetadata(container forge.Container) ContainerMetadata {
	m := ContainerMetadata{}
	m["id"] = container.GetId()
	return m
}

type ContainerMetadata map[string]string

func (m ContainerMetadata) GetId() string { //nolint:revive // matching protobuf style
	return m["id"]
}

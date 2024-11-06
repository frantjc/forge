package concourse

// Pipeline is the subset of a Concourse pipeline file relevant
// to finding resource_types.
type Pipeline struct {
	ResourceTypes []ResourceType `json:"resource_types,omitempty"`
	Resources     []Resource     `json:"resources,omitempty"`
}

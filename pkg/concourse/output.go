package concourse

type Output struct {
	Version  map[string]string `json:"version,omitempty"`
	Metadata []*OutputMetadata `json:"metadata,omitempty"`
}

type OutputMetadata struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

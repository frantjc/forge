package concourse

// Output has the JSON encoding which
// a Resource outputs on stdout.
type Output struct {
	Version  map[string]string `json:"version,omitempty"`
	Metadata []struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"metadata,omitempty"`
}

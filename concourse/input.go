package concourse

// Input is the struct which has the JSON encoding which
// gets passed on stdin to a resource.
type Input struct {
	Params  map[string]any `json:"params,omitempty"`
	Source  map[string]any `json:"source,omitempty"`
	Version map[string]any `json:"version,omitempty"`
}

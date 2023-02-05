package concourse

// Input is the struct which has the JSON encoding which
// gets passed on stdin to a Resource.
type Input struct {
	Params  map[string]string `json:"params,omitempty"`
	Source  map[string]string `json:"source,omitempty"`
	Version map[string]string `json:"version,omitempty"`
}

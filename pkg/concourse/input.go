package concourse

type Input struct {
	Params  map[string]string `json:"params,omitempty"`
	Source  map[string]string `json:"source,omitempty"`
	Version map[string]string `json:"version,omitempty"`
}

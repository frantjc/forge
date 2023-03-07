package circleci

type Config struct {
	Version float64           `json:"version,omitempty" yaml:",omitempty"`
	Orbs    map[string]string `json:"orbs,omitempty" yaml:",omitempty"`
}

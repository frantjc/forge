package circleci

// Config is the subset of a CircleCI configuration file relevant
// to finding Orb aliases.
type Config struct {
	Version float64           `json:"version,omitempty" yaml:",omitempty"`
	Orbs    map[string]string `json:"orbs,omitempty" yaml:",omitempty"`
}

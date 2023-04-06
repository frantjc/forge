package circleci

type Source struct {
	Version     float64 `json:"version,omitempty" yaml:",omitempty"`
	Description string  `json:"description,omitempty" yaml:",omitempty"`
	Display     *struct {
		HomeURL   string `json:"home_url,omitempty" yaml:",omitempty"`
		SourceURL string `json:"source_url,omitempty" yaml:",omitempty"`
	} `json:"display,omitempty"`
	Commands  map[string]Command `json:"commands,omitempty" yaml:",omitempty"`
	Executors map[string]struct {
		Description string `json:"description,omitempty" yaml:",omitempty"`
		Docker      []struct {
			Image string `json:"image,omitempty" yaml:",omitempty"`
		}
		Parameters Parameters `json:"parameters,omitempty" yaml:",omitempty"`
	} `json:"executors,omitempty"`
	Jobs map[string]Job
}

type Command struct {
	Description string     `json:"description,omitempty" yaml:",omitempty"`
	Parameters  Parameters `json:"parameters,omitempty" yaml:",omitempty"`
	Steps       []Step     `json:"steps,omitempty" yaml:",omitempty"`
}

type Job struct {
	Command  `json:",inline"`
	Executor map[string]any `json:"executor,omitempty" yaml:",omitempty"`
}

type Parameter struct {
	Default     any      `json:"default,omitempty" yaml:",omitempty"`
	Description string   `json:"description,omitempty" yaml:",omitempty"`
	Type        string   `json:"type,omitempty" yaml:",omitempty"`
	Enum        []string `json:"enum,omitempty" yaml:",omitempty"`
}

type Parameters map[string]Parameter

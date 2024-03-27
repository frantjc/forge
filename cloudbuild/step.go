package cloudbuild

// Step represents a CloudBuild step.
type Step struct {
	ID                   string            `json:"id,omitempty"`
	Name                 string            `json:"name,omitempty"`
	Entrypoint           string            `json:"entrypoint,omitempty"`
	Args                 []string          `json:"args,omitempty"`
	Script               string            `json:"script,omitempty"`
	Env                  []string          `json:"env,omitempty"`
	Substitutions        map[string]string `json:"substitutions,omitempty"`
	AutomapSubstitutions bool              `json:"automapSubstitutions,omitempty"`
	DynamicSubstitutions bool              `json:"dynamicSubstitutions,omitempty"`
}

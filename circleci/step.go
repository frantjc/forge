package circleci

// Step ...
//
//	steps:
//	  - when:
//	 	  condition:
//			equal:
//			  - 1.0.0
//			  - << parameters.version >>
//		  run:
//		    command: env
//	  - run:
//	      command: env
//	  - node/install:
//		  key: value
type Step struct {
	When   *Conditional `json:"when,omitempty" yaml:",omitempty"`
	Unless *Conditional `json:"unless,omitempty" yaml:",omitempty"`
	Run    *struct {
		Name             string            `json:"name,omitempty" yaml:",omitempty"`
		Command          string            `json:"command,omitempty" yaml:",omitempty"`
		WorkingDirectory string            `json:"working_directory,omitempty" yaml:",omitempty"`
		Environment      map[string]string `json:"environment,omitempty" yaml:",omitempty"`
	} `json:"run,omitempty" yaml:",omitempty"`
	Dynamic map[string]map[string]any `json:",inline" yaml:",inline"`
}

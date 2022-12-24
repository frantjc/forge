package githubactions

type Step struct {
	Shell      string            `json:"shell,omitempty"`
	If         string            `json:"if,omitempty"`
	Name       string            `json:"name,omitempty"`
	ID         string            `json:"id,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	WorkingDir string            `json:"working_dir,omitempty"`
	Uses       string            `json:"uses,omitempty"`
	With       map[string]string `json:"with,omitempty"`
	Run        string            `json:"run,omitempty"`
}

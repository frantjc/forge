package forgeactions

type Mapping struct {
	ActionPath       string `json:"action_path,omitempty"`
	Workspace        string `json:"workspace,omitempty"`
	RunnerToolCache  string `json:"runner_tool_cache,omitempty"`
	RunnerTemp       string `json:"runner_temp,omitempty"`
	GitHubPath       string `json:"git_hub_path,omitempty"`
	GitHubPathPath   string `json:"git_hub_path_path,omitempty"`
	GitHubEnvPath    string `json:"git_hub_env_path,omitempty"`
	GitHubOutputPath string `json:"git_hub_output_path,omitempty"`
	GitHubStatePath  string `json:"git_hub_state_path,omitempty"`
}

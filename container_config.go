package forge

// ContainerConfig is the configuration that is used to
// create a container or an exec in a running container.
type ContainerConfig struct {
	Entrypoint []string `json:"entrypoint,omitempty"`
	Cmd        []string `json:"cmd,omitempty"`
	WorkingDir string   `json:"working_dir,omitempty"`
	Env        []string `json:"env,omitempty"`
	User       string   `json:"user,omitempty"`
	Privileged bool     `json:"privileged,omitempty"`
	Mounts     []*Mount `json:"mounts,omitempty"`
}

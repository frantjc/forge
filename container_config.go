package forge

// ContainerConfig is the configuration that is used to
// create a container or an exec in a running container.
type ContainerConfig struct {
	Entrypoint []string
	Cmd        []string
	WorkingDir string
	Env        []string
	User       string
	Privileged bool
	Mounts     []Mount
}

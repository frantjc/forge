package forgeconcourse

import (
	"github.com/frantjc/forge/concourse"
)

type Config struct {
	ResourceTypes []concourse.ResourceType `yaml:"resource_types,omitempty"`
	Resources     []concourse.Resource     `yaml:"resources,omitempty"`
}

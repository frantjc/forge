package forgeconcourse

import (
	"github.com/frantjc/forge/pkg/concourse"
)

type Config struct {
	ResourceTypes []*concourse.ResourceType `yaml:"resource_types,omitempty"`
	Resources     []*concourse.Resource     `yaml:"resources,omitempty"`
}

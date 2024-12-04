package concourse

// ResourceType has the YAML encoding of a resource type as it would appear
// in the `resource_types` array of a Concourse pipeline configuration file.
type ResourceType struct {
	Name       string              `json:"name,omitempty"`
	Privileged bool                `json:"privileged,omitempty"`
	Source     *ResourceTypeSource `json:"source,omitempty"`
}

type ResourceTypeSource struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag,omitempty"`
}

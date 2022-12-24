package concourse

type Source struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `pjson:"tag,omitempty"`
}

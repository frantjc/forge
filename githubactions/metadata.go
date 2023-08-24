package githubactions

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

const (
	RunsUsingDockerfileImage   = "Dockerfile"
	RunsUsingDockerImagePrefix = "docker://"
	RunsUsingDocker            = "docker"
	RunsUsingComposite         = "composite"
	RunsUsingNode12            = "node12"
	RunsUsingNode16            = "node16"
)

var ErrMissingRequiredInput = errors.New("required input missing")

func NewMetadataFromReader(r io.Reader) (*Metadata, error) {
	m := &Metadata{}
	return m, yaml.NewDecoder(r).Decode(m)
}

func (m *Metadata) InputsFromWith(with map[string]string) (map[string]string, error) {
	inputs := make(map[string]string, len(m.Inputs))
	for name, input := range m.Inputs {
		w, ok := with[name]
		switch {
		case ok:
			inputs[name] = w
		case input.Default != "":
			inputs[name] = fmt.Sprint(input.Default)
		case input.Required:
			return nil, ErrMissingRequiredInput
		}
	}
	return inputs, nil
}

func (m *Metadata) IsComposite() bool {
	return m.Runs.Using == RunsUsingComposite
}

func (m *Metadata) IsDockerfile() bool {
	return m.Runs.Using == RunsUsingDocker && m.Runs.Image == RunsUsingDockerfileImage
}

type Metadata struct {
	Name        string                    `json:"name,omitempty"`
	Author      string                    `json:"author,omitempty"`
	Description string                    `json:"description,omitempty"`
	Inputs      map[string]MetadataInput  `json:"inputs,omitempty"`
	Output      map[string]MetadataOutput `json:"output,omitempty"`
	Runs        *MetadataRuns             `json:"runs,omitempty"`
}

type MetadataInput struct {
	Description        string `json:"description,omitempty"`
	Required           bool   `json:"required,omitempty"`
	Default            string `json:"default,omitempty"`
	DeprecationMessage string `json:"deprecation_message,omitempty"`
}

type MetadataOutput struct {
	Description string `json:"description,omitempty"`
}

type MetadataRuns struct {
	Plugin         string            `json:"plugin,omitempty"`
	Using          string            `json:"using,omitempty"`
	Pre            string            `json:"pre,omitempty"`
	Main           string            `json:"main,omitempty"`
	Post           string            `json:"post,omitempty"`
	Image          string            `json:"image,omitempty"`
	PreEntrypoint  string            `json:"pre_entrypoint,omitempty"`
	Entrypoint     string            `json:"entrypoint,omitempty"`
	PostEntrypoint string            `json:"post_entrypoint,omitempty"`
	Args           []string          `json:"args,omitempty"`
	Env            map[string]string `json:"env,omitempty"`
	Steps          []Step            `json:"steps,omitempty"`
}

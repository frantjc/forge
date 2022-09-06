package actions

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

const (
	RunsUsingDockerImagePrefix = "docker://"
	RunsUsingDocker            = "docker"
	RunsUsingComposite         = "composite"
	RunsUsingNode12            = "node12"
	RunsUsingNode16            = "node16"
)

var (
	ErrMissingRequiredInput = errors.New("required input missing")
)

func NewMetadataFromReader(r io.Reader) (*Metadata, error) {
	m := &Metadata{}
	d := yaml.NewDecoder(r)
	return m, d.Decode(m)
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

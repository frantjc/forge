package fc

import (
	"bytes"
	"encoding/json"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/concourse"
)

func NewStreams(drains *forge.Drains, input *concourse.Input) *forge.Streams {
	in := new(bytes.Buffer)

	if err := json.NewEncoder(in).Encode(input); err != nil {
		return &forge.Streams{
			Drains: drains,
		}
	}

	return &forge.Streams{
		In:     in,
		Drains: drains,
	}
}

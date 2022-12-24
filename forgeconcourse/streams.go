package forgeconcourse

import (
	"bytes"
	"encoding/json"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/concourse"
)

func NewStreams(drains *forge.Drains, input *concourse.Input) *forge.Streams {
	in := new(bytes.Buffer)

	if err := json.NewEncoder(in).Encode(input); err != nil {
		return drains.ToStreams(nil)
	}

	return drains.ToStreams(in)
}

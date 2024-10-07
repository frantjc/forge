package forgeconcourse

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/concourse"
	"github.com/frantjc/forge/internal/contaminate"
)

// NewStreams creates a Streams with the JSON encoding of input on stdin.
func NewStreams(ctx context.Context, drains *forge.Drains, input *concourse.Input) *forge.Streams {
	stdin := contaminate.StdinFrom(ctx)
	if stdin == nil {
		in := new(bytes.Buffer)

		if err := json.NewEncoder(in).Encode(input); err == nil {
			stdin = in
		}
	}

	return drains.ToStreams(stdin)
}

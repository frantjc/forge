package digestutil

import (
	"bytes"
	"encoding/json"

	"github.com/opencontainers/go-digest"
)

func JSON(a any) (digest.Digest, error) {
	r := new(bytes.Buffer)

	if err := json.NewEncoder(r).Encode(a); err != nil {
		return digest.Digest(""), err
	}

	return digest.FromReader(r)
}

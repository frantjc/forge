package forge

import (
	"encoding/json"

	"github.com/opencontainers/go-digest"
)

func Digest(a any) (digest.Digest, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}

	d := digest.FromBytes(b)
	return d, d.Validate()
}

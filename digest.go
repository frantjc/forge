package forge

import (
	"encoding/json"

	"github.com/opencontainers/go-digest"
)

// Digest returns a digest.Digest of anything.
// Used to get unique encodings of Ores to
// store and retrieve them by from a Basin.
func Digest(a any) (digest.Digest, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}

	d := digest.FromBytes(b)
	return d, d.Validate()
}

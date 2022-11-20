package docker

import (
	"compress/gzip"
	"encoding/json"
	"io"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/opencontainers/go-digest"
	imagespecsv1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type Image struct {
	v1.Image
	name.Reference
}

func (i *Image) Digest() (digest.Digest, error) {
	hash, err := i.Image.Digest()
	if err != nil {
		return "", err
	}

	d := digest.Digest(hash.String())
	return d, d.Validate()
}

func (i *Image) Manifest() (manifest *imagespecsv1.Manifest, _ error) {
	rawManifest, err := i.RawManifest()
	if err != nil {
		return nil, err
	}

	return manifest, json.Unmarshal(rawManifest, manifest)
}

func (i *Image) Config() (config *imagespecsv1.ImageConfig, _ error) {
	// RawConfigFile returns the JSON that has this structure:
	//
	// {
	// 	...
	// 	"config": { ... }
	// }
	//
	rawConfig, err := i.RawConfigFile()
	if err != nil {
		return nil, err
	}

	// We want "config" from the above JSON, so we create
	// this struct containing our ImageConfig where
	// the "config" will be unmarshaled to.
	return config, json.Unmarshal(rawConfig, &struct {
		Config *imagespecsv1.ImageConfig `json:"config"`
	}{
		Config: config,
	})
}

func (i *Image) MarshalJSON() ([]byte, error) {
	d, err := i.Digest()
	if err != nil {
		return nil, err
	}

	return []byte("\"" + d.String() + "\""), nil
}

func (i *Image) Blob() io.Reader {
	pr, pw := io.Pipe()

	go func() {
		w := gzip.NewWriter(pw)
		err := tarball.Write(i, i.Image, w)
		if err == nil {
			err = w.Close()
		}
		_ = pw.CloseWithError(err)
	}()

	return pr
}

func (i *Image) GoString() string {
	return "&Image{" + i.String() + "}"
}

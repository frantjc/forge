package dockerd

import (
	"compress/gzip"
	"encoding/json"
	"io"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	imagespecsv1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type Image struct {
	v1.Image
	name.Reference
}

func (i *Image) Config() (*imagespecsv1.ImageConfig, error) {
	// RawConfigFile returns JSON that has this structure:
	//
	// {
	// 	...
	// 	"config": { ... }
	// }
	//
	// We want "config" from the above JSON, so we create
	// this struct containing our ImageConfig where
	// the "config" will be unmarshaled to.
	configFile := &struct {
		Config *imagespecsv1.ImageConfig `json:"config"`
	}{}

	rawConfig, err := i.RawConfigFile()
	if err != nil {
		return nil, err
	}

	return configFile.Config, json.Unmarshal(rawConfig, configFile)
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

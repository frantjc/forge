package yaml

import (
	"io"

	"sigs.k8s.io/yaml"
)

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

type Decoder struct {
	r io.Reader
}

func (d *Decoder) Decode(obj any) error {
	b, err := io.ReadAll(d.r)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, obj)
}

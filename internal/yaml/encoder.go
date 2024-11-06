package yaml

import (
	"io"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/tools/clientcmd/api/latest"
	"sigs.k8s.io/yaml"
)

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

type Encoder struct {
	w io.Writer
}

func (e *Encoder) Encode(obj any) error {
	if robj, ok := obj.(runtime.Object); ok {
		if crobj, err := latest.Scheme.ConvertToVersion(robj, latest.ExternalVersion); err == nil {
			return new(printers.YAMLPrinter).PrintObj(crobj, e.w)
		}
	}

	b, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = e.w.Write(b)
	return err
}

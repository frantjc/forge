package bin

import (
	"bytes"
	"io"
)

func NewShimReader() io.Reader {
	return bytes.NewReader(shim)
}

package unixtable

import "bytes"

func Marshal(a any) ([]byte, error) {
	var (
		buf     = new(bytes.Buffer)
		encoder = NewEncoder(buf)
	)

	if err := encoder.Encode(a); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

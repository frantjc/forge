package tarutil

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"strings"
)

var ErrNoFilesWithPrefix = errors.New("no files in tarball have prefix")

func StripPrefix(r io.Reader, prefix string, opts ...Opt) io.ReadCloser {
	var (
		o                   = new(Opts)
		pr, pw              = io.Pipe()
		ir                  = r
		iw        io.Writer = pw
		found               = false
		lenPrefix           = len(prefix)
	)
	for _, opt := range opts {
		opt(o)
	}

	go func() {
		defer pw.Close()

		if o.gzip {
			zr, err := gzip.NewReader(r)
			if err != nil {
				_ = pw.CloseWithError(err)
				return
			}
			defer zr.Close()

			ir = zr
		}

		tr, tw := tar.NewReader(ir), tar.NewWriter(iw)
		defer tw.Close()

		for {
			f, err := tr.Next()
			if errors.Is(err, io.EOF) {
				if !found {
					_ = pw.CloseWithError(ErrNoFilesWithPrefix)
				}

				break
			} else if err != nil {
				_ = pw.CloseWithError(err)
				break
			}

			if !strings.HasPrefix(f.Name, prefix) {
				continue
			}

			found = true
			f.Name = f.Name[lenPrefix:]

			if f.Name == "" || f.Name == "/" {
				continue
			}

			if err := tw.WriteHeader(f); err != nil {
				_ = pw.CloseWithError(err)
				break
			}

			//nolint:gosec
			if _, err := io.Copy(tw, tr); err != nil {
				_ = pw.CloseWithError(err)
				break
			}
		}
	}()

	return pr
}

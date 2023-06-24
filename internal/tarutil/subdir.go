package tarutil

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"strings"
)

var ErrEmptySubdir = errors.New("empty tarball subdirectory")

// Subdir reads the tarball from r and and streams the files in the
// given subdirectory to the returned io.ReadCloser as a tarball with
// the subdirectory's path trimmed from each file's name.
//
// If the subdirectory is empty or non-existent, the returned io.ReadCloser
// is closed with ErrEmptySubdir.
func Subdir(r io.Reader, subdir string, opts ...Opt) io.ReadCloser {
	var (
		o                   = new(Opts)
		pr, pw              = io.Pipe()
		ir                  = r
		iw        io.Writer = pw
		found               = false
		lenSubdir           = len(subdir)
	)
	for _, opt := range opts {
		opt(o)
	}

	go func() {
		defer pw.Close()

		if o.gzipped {
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
					_ = pw.CloseWithError(ErrEmptySubdir)
				}

				break
			} else if err != nil {
				_ = pw.CloseWithError(err)
				break
			}

			if !strings.HasPrefix(f.Name, subdir) {
				continue
			}

			found = true
			f.Name = f.Name[lenSubdir:]

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

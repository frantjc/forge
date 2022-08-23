package bin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"time"
)

const (
	ShimName = "shim"
)

var (
	modTime = time.Date(1985, time.October, 26, 8, 15, 00, 0, time.UTC)
	size    = int64(len(shim))
	mode    = int64(0777)
)

func init() {
	_ = NewShimTarArchive()
}

func NewShimTarArchive() io.ReadCloser {
	rc, err := newShimTarArchive()
	if err != nil {
		panic(err)
	}

	return rc
}

func newShimTarArchive() (io.ReadCloser, error) {
	var (
		tarArchive      = new(bytes.Buffer)
		gzipWriter, err = gzip.NewWriterLevel(tarArchive, Compression)
		tarWriter       = tar.NewWriter(gzipWriter)
	)
	if err != nil {
		return nil, err
	}
	defer gzipWriter.Close()
	defer tarWriter.Close()

	if err := tarWriter.WriteHeader(&tar.Header{
		Name:    ShimName,
		Size:    size,
		Mode:    mode,
		ModTime: modTime,
	}); err != nil {
		return nil, err
	}

	if _, err := io.Copy(tarWriter, NewShimReader()); err != nil {
		return nil, err
	}

	return io.NopCloser(tarArchive), nil
}

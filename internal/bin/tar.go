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

func NewShimTarArchive() io.Reader {
	return bytes.NewReader(shimTarArchive)
}

var (
	modTime = time.Date(1985, time.October, 26, 8, 15, 00, 0, time.UTC)
	size    = int64(len(shim))
	mode    = int64(0777)
)

var (
	shimTarArchive []byte
)

func init() {
	var (
		tarArchiveBuf   = bytes.NewBuffer(shimTarArchive)
		gzipWriter, err = gzip.NewWriterLevel(tarArchiveBuf, TarArchiveCompression)
		tarWriter       = tar.NewWriter(gzipWriter)
	)
	if err != nil {
		panic(err)
	}
	defer gzipWriter.Close()
	defer tarWriter.Close()

	if err := tarWriter.WriteHeader(&tar.Header{
		Name:    ShimName,
		Size:    size,
		Mode:    mode,
		ModTime: modTime,
	}); err != nil {
		panic(err)
	}

	if _, err := io.Copy(tarWriter, NewShimReader()); err != nil {
		panic(err)
	}
}

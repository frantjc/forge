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
	modTime = time.Date(1985, time.October, 26, 8, 15, 0o0, 0, time.UTC)
	size    = int64(len(shim))
	mode    = int64(0o777)
)

func NewShimTarArchive() io.Reader {
	var (
		tarArchive      = new(bytes.Buffer)
		gzipWriter, err = gzip.NewWriterLevel(tarArchive, TarArchiveCompression)
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

	if _, err = io.Copy(tarWriter, NewShimReader()); err != nil {
		panic(err)
	}

	return tarArchive
}

func NewTarArchiveWithEmptyFiles(names ...string) io.Reader {
	var (
		tarArchive      = new(bytes.Buffer)
		gzipWriter, err = gzip.NewWriterLevel(tarArchive, TarArchiveCompression)
		tarWriter       = tar.NewWriter(gzipWriter)
	)
	if err != nil {
		panic(err)
	}
	defer gzipWriter.Close()
	defer tarWriter.Close()

	for _, name := range names {
		if err := tarWriter.WriteHeader(&tar.Header{
			Name:    name,
			Mode:    mode,
			ModTime: modTime,
		}); err != nil {
			panic(err)
		}
	}

	return tarArchive
}

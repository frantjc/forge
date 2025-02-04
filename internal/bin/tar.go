package bin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"time"
)

var TarArchiveCompression = gzip.BestCompression

var (
	modTime = time.Date(1985, time.October, 26, 8, 15, 0o0, 0, time.UTC)
	size    = int64(len(shim))
	mode    = int64(0o777)
)

func NewShimReader() io.Reader {
	return bytes.NewReader(shim)
}

func NewShimTarArchive(path string) io.Reader {
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
		Name:    path,
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

func NewScriptTarArchive(script, path string) io.Reader {
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
		Name:    path,
		Size:    int64(len(script)),
		Mode:    mode,
		ModTime: modTime,
	}); err != nil {
		panic(err)
	}

	if _, err = io.Copy(tarWriter, strings.NewReader(script)); err != nil {
		panic(err)
	}

	return tarArchive
}

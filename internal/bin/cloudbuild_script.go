package bin

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"path/filepath"
	"strings"

	"github.com/frantjc/forge/internal/containerfs"
)

var (
	ScriptPath       = filepath.Join(containerfs.WorkingDir, ScriptName)
	ScriptEntrypoint = []string{ScriptPath}
)

const (
	ScriptName = "script"
)

func NewScriptTarArchive(script string) io.Reader {
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
		Name:    ScriptName,
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

package docker

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"

	xos "github.com/frantjc/x/os"
	imagespecsv1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type Image struct {
	Ref  string
	Path string
}

func (i *Image) Name() string { return i.Ref }

func (i *Image) Config() (*imagespecsv1.ImageConfig, error) {
	var (
		configFile = &struct {
			Config *imagespecsv1.ImageConfig `json:"config"`
		}{}
		buf = new(bytes.Buffer)
		//nolint:gosec
		cmd = exec.Command(i.Path, "inspect", i.Ref)
	)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, xos.NewExitCodeError(err, cmd.ProcessState.ExitCode())
	}

	return configFile.Config, json.NewDecoder(buf).Decode(configFile)
}

func (i *Image) Blob() io.Reader {
	pr, pw := io.Pipe()

	go func() {
		//nolint:gosec
		cmd := exec.Command(i.Path, "save", i.Ref)
		cmd.Stdout = pw
		if err := cmd.Run(); err != nil {
			_ = pw.CloseWithError(xos.NewExitCodeError(err, cmd.ProcessState.ExitCode()))
		} else {
			_ = pw.Close()
		}
	}()

	return pr
}

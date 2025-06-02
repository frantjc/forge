package docker

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
	"strings"

	xos "github.com/frantjc/x/os"
	"github.com/opencontainers/go-digest"
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
		cmd = exec.Command(i.Path, "inspect", i.Ref)
	)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, xos.NewExitCodeError(err, cmd.ProcessState.ExitCode())
	}

	return configFile.Config, json.NewDecoder(buf).Decode(configFile)
}

func (i *Image) Digest() (digest.Digest, error) {
	cmd := exec.Command(i.Path, "inspect", "--format={{.Id}}", i.Ref)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return digest.FromString(strings.TrimSpace(string(out))), nil
}

func (i *Image) Blob() io.Reader {
	pr, pw := io.Pipe()

	go func() {
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

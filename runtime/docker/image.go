package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		buf = new(bytes.Buffer)
		//nolint:gosec
		cmd = exec.Command(i.Path, "inspect", i.Ref)
	)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, xos.NewExitCodeError(err, cmd.ProcessState.ExitCode())
	}

	var cfgs []struct {
		Config *imagespecsv1.ImageConfig `json:"config"`
	}
	if err := json.NewDecoder(buf).Decode(&cfgs); err != nil {
		return nil, err
	}
	if len(cfgs) == 0 {
		return nil, fmt.Errorf("no inspect results for %s", i.Ref)
	}

	return cfgs[0].Config, nil
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

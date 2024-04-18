package forgecloudbuild

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/bin"
)

func CopyScriptToContainer(ctx context.Context, container forge.Container, script string) error {
	return DefaultMapping.CopyScriptToContainer(ctx, container, script)
}

func (m *Mapping) CopyScriptToContainer(ctx context.Context, container forge.Container, script string) error {
	_script := script
	if !bin.HasShebang(_script) {
		_script = fmt.Sprintf("#!/bin/sh\n%s", _script)
	}

	if err := container.CopyTo(ctx, filepath.Dir(bin.ScriptPath), bin.NewScriptTarArchive(_script)); err != nil {
		return err
	}

	return nil
}

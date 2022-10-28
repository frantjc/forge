package actions

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/frantjc/forge"
)

func (u *Uses) IsLocal() bool {
	return strings.HasPrefix(u.Path, "./") || filepath.IsAbs(u.Path) || len(strings.Split(u.Path, "/")) < 2
}

func (u *Uses) IsRemote() bool {
	return !u.IsLocal()
}

func (u *Uses) Uses() string {
	return u.GetPath() + "@" + u.GetVersion()
}

func (u *Uses) Repository() string {
	if u.IsRemote() {
		return strings.Join(strings.Split(u.Path, "/")[0:2], "/")
	}

	return ""
}

func (u *Uses) GoString() string {
	return "&Uses{" + u.Uses() + "}"
}

func Parse(uses string) (*Uses, error) {
	r := &Uses{}

	spl := strings.Split(uses, "@")
	switch len(spl) {
	case 2:
		r.Version = spl[1]
		fallthrough
	case 1:
		r.Path = spl[0]
	default:
		return r, fmt.Errorf("parse uses: '%s'", uses)
	}

	return r, nil
}

func (u *Uses) MarshalJSON() ([]byte, error) {
	return []byte("\"" + u.Uses() + "\""), nil
}

func GetUsesMetadata(ctx context.Context, uses *Uses, dir string) (*Metadata, error) {
	var (
		_ = forge.LoggerFrom(ctx)
		u = GetGitHubURL()
	)

	if uses.IsRemote() {
		return CloneUses(ctx, uses, &CloneOpts{
			GitHubURL: u,
			Insecure:  u.Scheme != "https",
			Path:      path.Clean(dir),
		})
	}

	r, err := OpenDirectoryMetadata(filepath.Clean(uses.Path))
	if err != nil {
		return nil, err
	}

	return NewMetadataFromReader(r)
}

func OpenUsesMetadata(uses *Uses) (io.Reader, error) {
	if uses.IsRemote() {
		return nil, fmt.Errorf("open remote action: %s", uses.Path)
	}

	return OpenDirectoryMetadata(filepath.Clean(uses.Path))
}

func OpenDirectoryMetadata(dir string) (_ io.Reader, err error) {
	for _, filename := range ActionYAMLFilenames {
		if f, err := os.Open(filepath.Join(dir, filename)); err == nil {
			return f, nil
		}
	}

	return nil, err
}

package githubactions

import (
	"context"
	"fmt"
	"io"
	"os"
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

func (u *Uses) UsesString() string {
	uses := u.GetPath()
	if v := u.GetVersion(); v != "" {
		uses = uses + "@" + v
	}
	return uses
}

func (u *Uses) GetRepository() string {
	if u.IsRemote() {
		return filepath.Join(strings.Split(u.GetPath(), "/")[0:2]...)
	}

	return ""
}

func (u *Uses) GetActionPath() string {
	if u.IsRemote() {
		elements := strings.Split(u.GetPath(), "/")
		if len(elements) > 2 {
			return filepath.Join(elements[2:]...)
		}
	}

	return ""
}

func (u *Uses) GoString() string {
	return "&Uses{" + u.UsesString() + "}"
}

// TODO regexp.
func Parse(uses string) (*Uses, error) {
	r := &Uses{}

	switch {
	case strings.HasPrefix(uses, "/"):
		r.Path = filepath.Clean(uses)
	case strings.HasPrefix(uses, "./"), strings.HasPrefix(uses, "../"):
		r.Path = "./" + filepath.Clean(uses)
	default:
		spl := strings.Split(uses, "@")
		if len(spl) != 2 {
			return nil, fmt.Errorf("parse uses: not a path or a versioned reference: %s", uses)
		}

		r.Path = filepath.Clean(spl[0])
		r.Version = spl[1]
	}

	return r, nil
}

func (u *Uses) MarshalJSON() ([]byte, error) {
	return []byte("\"" + u.UsesString() + "\""), nil
}

func GetUsesMetadata(ctx context.Context, uses *Uses, dir string) (*Metadata, error) {
	var (
		_ = forge.LoggerFrom(ctx)
		u = GetGitHubURL()
	)

	if uses.IsRemote() {
		return CheckoutUses(ctx, uses, &CheckoutOpts{
			GitHubURL: u,
			Insecure:  u.Scheme != "https",
			Path:      filepath.Clean(dir),
		})
	}

	r, err := OpenDirectoryMetadata(filepath.Join(dir, uses.GetActionPath()))
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

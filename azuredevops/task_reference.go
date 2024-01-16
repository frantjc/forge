package azuredevops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type TaskReference struct {
	Path    string
	Version string
}

func (r *TaskReference) IsLocal() bool {
	return strings.HasPrefix(r.Path, ".") || filepath.IsAbs(r.Path) || len(strings.Split(r.Path, "/")) < 2
}

func (r *TaskReference) IsRemote() bool {
	return !r.IsLocal()
}

func (r *TaskReference) String() string {
	ref := r.Path
	if v := r.Version; v != "" {
		ref = ref + "@" + v
	}
	return ref
}

// TODO regexp.
func Parse(ref string) (*TaskReference, error) {
	r := &TaskReference{}

	switch {
	case strings.HasPrefix(ref, "/"):
		r.Path = filepath.Clean(ref)
	case strings.HasPrefix(ref, "."):
		r.Path = filepath.Clean(ref)
	default:
		spl := strings.Split(ref, "@")
		if len(spl) != 2 {
			return nil, fmt.Errorf("parse task reference: not a path or a versioned reference: %s", ref)
		}

		r.Path = filepath.Clean(spl[0])
		r.Version = spl[1]
	}

	if r.Path != "." && !filepath.IsAbs(r.Path) {
		r.Path = "./" + r.Path
	}

	return r, nil
}

func (r *TaskReference) MarshalJSON() ([]byte, error) {
	return []byte("\"" + r.String() + "\""), nil
}

func GetReferenceTask(ref *TaskReference) (*Task, error) {
	r, err := OpenReferenceTask(ref)
	if err != nil {
		return nil, err
	}

	return NewTaskFromReader(r)
}

func OpenReferenceTask(ref *TaskReference) (io.Reader, error) {
	if ref.IsRemote() {
		return nil, fmt.Errorf("open remote task: %s", ref.Path)
	}

	return OpenDirectoryTask(filepath.Clean(ref.Path))
}

func OpenDirectoryTask(dir string) (io.Reader, error) {
	name, err := filepath.Abs(filepath.Join(dir, "task.json"))
	if err != nil {
		return nil, err
	}

	return os.Open(name)
}

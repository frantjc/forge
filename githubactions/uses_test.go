package githubactions_test

import (
	"testing"

	"github.com/frantjc/forge/githubactions"
)

func TestParse(t *testing.T) {
	for _, s := range []struct {
		uses            string
		expectedPath    string
		expectedVersion string
		expectedLocal   bool
	}{
		{
			uses:            "./frantjc/forge",
			expectedPath:    "./frantjc/forge",
			expectedVersion: "",
			expectedLocal:   true,
		},
		{
			uses:            ".",
			expectedPath:    ".",
			expectedVersion: "",
			expectedLocal:   true,
		},
		{
			uses:            "frantjc/forge@v0",
			expectedPath:    "frantjc/forge",
			expectedVersion: "v0",
		},
	} {
		if actual, err := githubactions.Parse(s.uses); err != nil {
			t.Error(err)
			t.FailNow()
		} else if actual.Path != s.expectedPath {
			t.Error("was", actual.Path, "but expected", s.expectedPath)
			t.FailNow()
		} else if actual.Version != s.expectedVersion {
			t.Error("was", actual.Version, "but expected", s.expectedVersion)
			t.FailNow()
		} else if local := actual.IsLocal(); local != s.expectedLocal {
			t.Error("was", local, "but expected", s.expectedLocal)
			t.FailNow()
		}
	}
}

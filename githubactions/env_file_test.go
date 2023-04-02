package githubactions_test

import (
	"bytes"
	"testing"

	"github.com/frantjc/forge/githubactions"
	"github.com/google/uuid"
)

func TestParseEnvFile(t *testing.T) {
	var (
		env = `# comment
HELLO=there

GENERAL<<ghadelimiter_` + uuid.NewString() + `
kenobi # comment
ghadelimiter_` + uuid.NewString() + `

YOU="are a"

BOLD<<ghadelimiter_` + uuid.NewString() + `
one
ghadelimiter_` + uuid.NewString() + `
		`
		expected = map[string]string{
			"HELLO":   "there",
			"GENERAL": "kenobi",
			"YOU":     "are a",
			"BOLD":    "one",
		}
		actual, err = githubactions.ParseEnvFile(bytes.NewBufferString(env))
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for k, v := range actual {
		if expected[k] != v {
			t.Error("actual", v, "for key", k, "does not match expected", expected[k])
			t.FailNow()
		}
	}
}

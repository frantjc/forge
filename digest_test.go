package forge_test

import (
	"testing"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/pkg/github/actions"
	"github.com/frantjc/forge/pkg/ore"
)

func TestDigest(t *testing.T) {
	var (
		o = &ore.Action{
			Uses:          "actions/checkout@v3",
			GlobalContext: actions.NewGlobalContextFromEnv(),
		}
	)

	expected, err := forge.Digest(o)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	actual, err := forge.Digest(o)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if actual != expected {
		t.Error("actual digest", actual.Encoded(), "is not equal to expected digest", expected.Encoded())
		t.FailNow()
	}
}

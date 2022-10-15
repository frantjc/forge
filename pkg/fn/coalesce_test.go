package fn_test

import (
	"testing"

	"github.com/frantjc/forge/pkg/fn"
)

func TestCoalesceTrue(t *testing.T) {
	var (
		expected = "default"
		actual   = fn.Coalesce("", "", "", expected, "")
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestCoalesceFalse(t *testing.T) {
	var (
		expected = ""
		actual   = fn.Coalesce("", "", "")
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

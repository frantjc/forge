package fn_test

import (
	"testing"

	"github.com/frantjc/forge/pkg/fn"
)

func TestTernaryTrue(t *testing.T) {
	var (
		expected = 1
		actual   = fn.Ternary(true, 1, 0)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestTernaryFalse(t *testing.T) {
	var (
		expected = 0
		actual   = fn.Ternary(false, 1, 0)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

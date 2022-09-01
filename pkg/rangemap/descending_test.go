package rangemap_test

import (
	"testing"

	"github.com/frantjc/forge/pkg/rangemap"
)

func TestDescending(t *testing.T) {
	var (
		m = map[string]int{
			"b": 3,
			"d": 1,
			"c": 2,
			"a": 4,
		}
		i = 1
	)

	rangemap.Descending(m, func(k string, v int) {
		if i != v {
			t.Error(k, "was in position", v, "should have been in position", i)
			t.FailNow()
		}

		i++
	})
}

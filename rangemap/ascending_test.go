package rangemap_test

import (
	"testing"

	"github.com/frantjc/forge/rangemap"
)

func TestAscending(t *testing.T) {
	var (
		m = map[string]int{
			"b": 2,
			"d": 4,
			"c": 3,
			"a": 1,
		}
		i = 1
	)

	rangemap.Ascending(m, func(k string, v int) {
		if i != v {
			t.Error(k, "was in position", v, "should have been in position", i)
			t.FailNow()
		}

		i++
	})
}

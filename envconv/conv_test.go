package envconv_test

import (
	"testing"

	"github.com/frantjc/forge/envconv"
)

func TestMapFromArr(t *testing.T) {
	var (
		a        = []string{"KEY1=val", "KEY2=", "=val", "notakeyvalpair"}
		expected = map[string]string{
			"KEY1": "val",
			"KEY2": "",
		}
		actual = envconv.ArrToMap(a)
	)

	for k, v := range actual {
		if expected[k] != v {
			t.Error("key", k, "actual value", v, "does not equal expected value", expected[k])
			t.FailNow()
		}
	}
}

func TestToMap(t *testing.T) {
	var (
		expected = map[string]string{
			"KEY1": "val",
			"KEY2": "",
		}
		actual = envconv.ToMap("KEY1=val", "KEY2=", "=val", "notakeyvalpair")
	)

	for k, v := range actual {
		if expected[k] != v {
			t.Error("key", k, "actual value", v, "does not equal expected value", expected[k])
			t.FailNow()
		}
	}
}

func TestArrFromMap(t *testing.T) {
	var (
		m = map[string]string{
			"KEY1": "val",
			"":     "val",
			"KEY2": "",
		}
		expected = []string{"KEY1=val", "KEY2="}
		actual   = envconv.MapToArr(m)
	)

	for _, a := range actual {
		contains := false

		for _, e := range expected {
			if a == e {
				contains = true
			}
		}

		if !contains {
			t.Error("actual contains", a, "but expected", expected, "does not")
			t.FailNow()
		}
	}
}

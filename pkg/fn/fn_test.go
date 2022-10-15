package fn_test

import (
	"testing"

	"github.com/frantjc/forge/pkg/fn"
)

func TestEveryTrue(t *testing.T) {
	var (
		a = []int{1, 2, 3, 4}
		f = func(a, _ int) bool {
			return a > 0
		}
		expected = true
		actual   = fn.Every(a, f)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestEveryFalse(t *testing.T) {
	var (
		a = []int{0, 1, 2, 3}
		f = func(a, _ int) bool {
			return a > 0
		}
		expected = false
		actual   = fn.Every(a, f)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestFilter(t *testing.T) {
	var (
		a = []int{0, 1, 2, 3}
		f = func(a, _ int) bool {
			return a > 0
		}
		expected = []int{1, 2, 3}
		actual   = fn.Filter(a, f)
	)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Error("actual", actual, "does not equal expected", expected)
			t.FailNow()
		}
	}
}

func TestFindIndex(t *testing.T) {
	var (
		a = []int{1, 2, 3, 4}
		f = func(a, _ int) bool {
			return a == 2
		}
		expected = 1
		actual   = fn.FindIndex(a, f)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestFind(t *testing.T) {
	var (
		a = []int{1, 2, 3, 4}
		f = func(a, _ int) bool {
			return a > 0
		}
		expected = 1
		actual   = fn.Find(a, f)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestIncludesTrue(t *testing.T) {
	var (
		a        = []int{1, 2, 3, 4}
		expected = true
		actual   = fn.Includes(a, 1)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestIncludesFalse(t *testing.T) {
	var (
		a        = []int{1, 2, 3, 4}
		expected = false
		actual   = fn.Includes(a, 0)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestIndexOf(t *testing.T) {
	var (
		a        = []int{1, 2, 3, 4}
		expected = 0
		actual   = fn.IndexOf(a, 1)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestLastIndexOf(t *testing.T) {
	var (
		a        = []int{1, 2, 3, 4, 3, 2}
		expected = 4
		actual   = fn.LastIndexOf(a, 3, 0)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestLastIndexOfFrom(t *testing.T) {
	var (
		a        = []int{1, 2, 3, 4, 3, 2}
		expected = 2
		actual   = fn.LastIndexOf(a, 3, 3)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestMap(t *testing.T) {
	var (
		a = []int{1, 2, 3, 4, 3, 2}
		f = func(a, _ int) int {
			return a + 1
		}
		expected = []int{2, 3, 4, 5, 4, 3}
		actual   = fn.Map(a, f)
	)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Error("actual", actual, "does not equal expected", expected)
			t.FailNow()
		}
	}
}

func TestReduceRight(t *testing.T) {
	var (
		a = []int{1, 2, 3, 4, 3, 2}
		f = func(acc, a, _ int) int {
			return acc + a
		}
		expected = 15
		actual   = fn.ReduceRight(a, f, 0)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestReduce(t *testing.T) {
	var (
		a = []int{1, 2, 3, 4, 3, 2}
		f = func(acc, a, _ int) int {
			return acc + a
		}
		expected = 15
		actual   = fn.Reduce(a, f, 0)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestReverse(t *testing.T) {
	var (
		a        = []int{1, 2, 3, 4, 3, 2}
		expected = []int{2, 3, 4, 3, 2, 1}
		actual   = fn.Reverse(a)
	)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Error("actual", actual, "does not equal expected", expected)
			t.FailNow()
		}
	}
}

func TestSliceFromStart(t *testing.T) {
	var (
		a        = []int{0, 1, 2, 3}
		expected = []int{0, 1}
		actual   = fn.Slice(a, 0, 2)
	)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Error("actual", actual, "does not equal expected", expected)
			t.FailNow()
		}
	}
}

func TestSliceFromEnd(t *testing.T) {
	var (
		a        = []int{0, 1, 2, 3}
		expected = []int{2}
		actual   = fn.Slice(a, -2, -1)
	)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Error("actual", actual, "does not equal expected", expected)
			t.FailNow()
		}
	}
}

func TestSomeTrue(t *testing.T) {
	var (
		a = []int{1, 2, 3, 4}
		f = func(a, _ int) bool {
			return a == 4
		}
		expected = true
		actual   = fn.Some(a, f)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestSomeFalse(t *testing.T) {
	var (
		a = []int{1, 2, 3, 4}
		f = func(a, _ int) bool {
			return a == 5
		}
		expected = false
		actual   = fn.Some(a, f)
	)

	if expected != actual {
		t.Error("actual", actual, "does not equal expected", expected)
		t.FailNow()
	}
}

func TestUnique(t *testing.T) {
	var (
		a        = []int{1, 2, 3, 4, 3, 2, 1}
		expected = []int{1, 2, 3, 4}
		actual   = fn.Unique(a)
	)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Error("actual", actual, "does not equal expected", expected)
			t.FailNow()
		}
	}
}

func TestSortWorst(t *testing.T) {
	var (
		a        = []int{4, 3, 2, 1}
		expected = []int{1, 2, 3, 4}
		actual   = fn.Sort(a, func(a, b int) int {
			return a - b
		})
	)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Error("actual", actual, "does not equal expected", expected)
			t.FailNow()
		}
	}
}

func TestSortBest(t *testing.T) {
	var (
		a        = []int{1, 2, 3, 4}
		expected = []int{1, 2, 3, 4}
		actual   = fn.Sort(a, func(a, b int) int {
			return a - b
		})
	)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Error("actual", actual, "does not equal expected", expected)
			t.FailNow()
		}
	}
}

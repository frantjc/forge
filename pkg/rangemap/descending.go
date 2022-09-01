package rangemap

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func Descending[T1 constraints.Ordered, T2 any](m map[T1]T2, f func(T1, T2)) {
	var (
		keys = make([]T1, len(m))
		i    = 0
	)

	for k := range m {
		keys[i] = k
		i++
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	for _, k := range keys {
		f(k, m[k])
	}
}

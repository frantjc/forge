package rangemap

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func ordered[T1 constraints.Ordered, T2 any](m map[T1]T2, f func(T1, T2), asc bool) {
	var (
		keys  = make([]T1, len(m))
		keysC = make(chan any, 1)
		i     = 0
	)
	go func() {
		for k := range m {
			keys[i] = k
			i++
		}
		close(keysC)
	}()

	var (
		less = func(i, j int) bool {
			return keys[i] > keys[j]
		}
	)
	if asc {
		less = func(i, j int) bool {
			return keys[i] < keys[j]
		}
	}

	<-keysC

	sort.Slice(keys, less)

	for _, k := range keys {
		f(k, m[k])
	}
}

package rangemap

import "golang.org/x/exp/constraints"

// Ascending takes a map with orderable keys and invokes
// the given callback function with the map's key-value pairs
// in ascending order by key.
func Ascending[T1 constraints.Ordered, T2 any](m map[T1]T2, f func(T1, T2)) {
	ordered(m, f, true)
}

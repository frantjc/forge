package rangemap

import "golang.org/x/exp/constraints"

// Descending takes a map with orderable keys and invokes
// the given callback function with the map's key-value pairs
// in descending order by key.
func Descending[T1 constraints.Ordered, T2 any](m map[T1]T2, f func(T1, T2)) {
	ordered(m, f, false)
}

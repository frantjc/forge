package rangemap

import "golang.org/x/exp/constraints"

func Ascending[T1 constraints.Ordered, T2 any](m map[T1]T2, f func(T1, T2)) {
	ordered(m, f, true)
}

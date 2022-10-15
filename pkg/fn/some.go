package fn

// Some tests whether at least one element in the array passes
// the test implemented by the provided function.
// It returns true if, in the array, it finds an element
// for which the provided function returns true; otherwise
// it returns false. It doesn't modify the array.
func Some[T any](in []T, f func(T, int) bool) bool {
	for i, a := range in {
		if f(a, i) {
			return true
		}
	}
	return false
}

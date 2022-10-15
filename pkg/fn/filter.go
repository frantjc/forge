package fn

// Filter creates a new array with all elements that pass
// the test implemented by the provided function.
func Filter[T any](in []T, f func(T, int) bool) []T {
	out := []T{}
	for i, a := range in {
		if f(a, i) {
			out = append(out, a)
		}
	}
	return out
}

package fn

// Find returns the first element in the provided array that satisfies
// the provided testing function. If no values satisfy the testing function,
// the zero value of T is returned.
func Find[T any](in []T, f func(T, int) bool) T {
	for i, a := range in {
		if f(a, i) {
			return a
		}
	}
	return *new(T) //nolint:gocritic
}

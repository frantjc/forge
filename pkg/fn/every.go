package fn

// Every tests whether all elements in the array pass the
// test implemented by the provided function. It returns a Boolean value.
func Every[T any](in []T, f func(T, int) bool) bool {
	for i, a := range in {
		if !f(a, i) {
			return false
		}
	}
	return true
}

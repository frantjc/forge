package fn

// Unique creates a new array with only unique elements from the given array.
func Unique[T comparable](in []T) []T {
	return Filter(in, func(a T, i int) bool {
		return i == IndexOf(in, a)
	})
}

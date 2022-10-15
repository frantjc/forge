package fn

// IndexOf returns the first index at which a given element
// can be found in the array, or -1 if it is not present.
func IndexOf[T comparable](in []T, a T) int {
	return FindIndex(in, func(b T, _ int) bool {
		return a == b
	})
}

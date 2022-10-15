package fn

// LastIndexOf returns the last index at which a given element
// can be found in the array, or -1 if it is not present.
// The array is searched backwards, starting at fromIndex.
func LastIndexOf[T comparable](in []T, a T, from int) int {
	if from <= 0 {
		from = len(in) - 1
	}
	for i := from; i >= 0; i-- {
		if a == in[i] {
			return i
		}
	}
	return -1
}

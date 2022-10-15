package fn

// Reverse creates a new array by reversing the given array.
// The first array element becomes the last, and the last array element becomes the first.
func Reverse[T any](in []T) []T {
	j := len(in)
	out := make([]T, j)
	k := j - 1
	for i, a := range in {
		out[k-i] = a
	}
	return out
}

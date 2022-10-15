package fn

// Sort creates a new array by sorting the elements of
// the given array and returns the sorted array.
// The sort order is ascending.
func Sort[T comparable](in []T, f func(a, b T) int) []T {
	out := in
	k := len(out)
	for i := 0; i < k; i++ {
		for j := 0; j < k-1; j++ {
			if f(out[j], out[j+1]) > 0 {
				out[j], out[j+1] = out[j+1], out[j]
			}
		}
	}
	return out
}

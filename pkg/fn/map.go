package fn

// Map creates a new array populated with the results of
// calling a provided function on every element in the calling array.
func Map[T1, T2 any](in []T1, f func(T1, int) T2) []T2 {
	out := make([]T2, len(in))
	for i, a := range in {
		out[i] = f(a, i)
	}
	return out
}

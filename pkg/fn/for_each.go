package fn

// ForEach executes a provided function once
// for each array element.
func ForEach[T any](in []T, f func(T, int)) {
	for i, a := range in {
		f(a, i)
	}
}

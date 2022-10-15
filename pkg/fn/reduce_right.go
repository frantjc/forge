package fn

// ReduceRight applies a function against an accumulator and each value of the array
// (from right-to-left) to reduce it to a single value.
func ReduceRight[T1, T2 any](in []T1, f func(T2, T1, int) T2, initial T2) T2 {
	accumulator := initial
	for i := len(in) - 1; i >= 0; i-- {
		accumulator = f(accumulator, in[i], i)
	}
	return accumulator
}

package fn

// Ternary returns the first T if the conditional evaluates to true,
// otherwise it returns the second T.
func Ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

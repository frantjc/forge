package fn

// Includes determines whether an array includes a certain value
// among its entries, returning true or false as appropriate.
func Includes[T comparable](in []T, a T) bool {
	return Some(in, func(b T, _ int) bool {
		return a == b
	})
}

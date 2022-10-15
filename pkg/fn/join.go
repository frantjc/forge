package fn

import "fmt"

// Join creates and returns a new string by concatenating all of the elements
// in an array (or an array-like object), separated by commas or a specified
// separator string. If the array has only one item, then that item will be
// returned without using the separator.
func Join[T fmt.Stringer](in []T, separator string) string {
	out := ""
	i := len(in)
	for j, a := range in {
		if 0 < j && j < i {
			out += separator
		}
		out += a.String()
	}
	return out
}

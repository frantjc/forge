package rangemap

import "sort"

func Alphabetically[T any](m map[string]T, f func(string, T)) {
	var (
		keys = make([]string, len(m))
		i    = 0
	)

	for k := range m {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for _, k := range keys {
		f(k, m[k])
	}
}

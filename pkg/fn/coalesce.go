package fn

func Coalesce[T comparable](ts ...T) T {
	empty := *new(T) //nolint:gocritic
	for _, t := range ts {
		if t != empty {
			return t
		}
	}
	return empty
}

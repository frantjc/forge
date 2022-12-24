package changroup

// AllSettled takes an arbitrary number of channels and returns
// a single channel that will be closed once all of the given
// channels have either been closed or received one value through
// them. Inspired by JavaScript's Promise.allSettled.
func AllSettled[T any](chans ...chan T) chan []T {
	var (
		t        = make([]T, len(chans))
		combined = make(chan []T, 1)
	)

	go func() {
		defer func() {
			if len(t) > 0 {
				combined <- t
			}
			close(combined)
		}()

		for _, c := range chans {
			t = append(t, <-c)
		}
	}()

	return combined
}

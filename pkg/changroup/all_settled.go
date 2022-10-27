package changroup

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

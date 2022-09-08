package channels

import "context"

func Wait(ctx context.Context, chans ...chan any) chan error {
	combined := make(chan error, 1)
	go func() {
		for _, c := range chans {
			select {
			case <-c:
			case <-ctx.Done():
				combined <- ctx.Err()
			}
		}
		close(combined)
	}()
	return combined
}

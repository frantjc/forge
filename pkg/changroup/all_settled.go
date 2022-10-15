package changroup

import "context"

type SettledResult[T any] struct {
	Err   error
	Value T
}

type AllSettledResult[T any] []*SettledResult[T]

func (a AllSettledResult[T]) Err() error {
	for _, r := range a {
		if r.Err != nil {
			return r.Err
		}
	}

	return nil
}

func (a AllSettledResult[T]) Values() []T {
	values := []T{}

	for _, r := range a {
		if r.Err != nil {
			values = append(values, r.Value)
		}
	}

	return values
}

func AllSettled[T any](ctx context.Context, chans ...chan T) chan AllSettledResult[T] {
	var (
		allSettledResult = AllSettledResult[T]{}
		combined         = make(chan AllSettledResult[T], 1)
	)

	go func() {
		defer func() {
			if len(allSettledResult) > 0 {
				combined <- allSettledResult
			}
			close(combined)
		}()

		for _, c := range chans {
			select {
			case result := <-c:
				allSettledResult = append(allSettledResult, &SettledResult[T]{Value: result})
			case <-ctx.Done():
				allSettledResult = append(allSettledResult, &SettledResult[T]{Err: ctx.Err()})
				return
			}
		}
	}()

	return combined
}

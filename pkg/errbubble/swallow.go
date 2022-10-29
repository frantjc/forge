package errbubble

import "errors"

func Swallow(err error, target error) error {
	if errors.Is(err, target) {
		return nil
	}

	return err
}

package errbubble

import "net/http"

type ExitCoder interface {
	ExitCode() int
}

func ExitCodeOr(err error, fallback int) int {
	if err != nil {
		if ec, ok := err.(ExitCoder); ok {
			return ec.ExitCode()
		}

		return 1
	}

	return fallback
}

func ExitCode(err error) int {
	return ExitCodeOr(err, 0)
}

type HTTPStatusCoder interface {
	HTTPStatusCode() int
}

func HTTPStatusCodeOr(err error, fallback int) int {
	if err != nil {
		if hsc, ok := err.(HTTPStatusCoder); ok {
			return hsc.HTTPStatusCode()
		}

		return http.StatusInternalServerError
	}

	return fallback
}

func HTTPStatusCode(err error) int {
	return HTTPStatusCodeOr(err, 0)
}

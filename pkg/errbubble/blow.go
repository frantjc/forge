package errbubble

type ExitCoder interface {
	ExitCode() int
}

func ExitCode(err error) int {
	return New(err).ExitCode()
}

type HTTPStatusCoder interface {
	HTTPStatusCode() int
}

func HTTPStatusCode(err error) int {
	return New(err).HTTPStatusCode()
}

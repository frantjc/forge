package errbubble

type Bubble struct {
	error
	exitCode, httpStatusCode int
}

func (b *Bubble) ExitCode() int {
	return b.exitCode
}

func (b *Bubble) HTTPStatusCode() int {
	return b.httpStatusCode
}

func New(err error, exitCode, statusCode int) error {
	return &Bubble{err, exitCode, statusCode}
}

func NewExitCode(err error, exitCode int) error {
	return New(err, exitCode, 0)
}

func NewHTTPStatusCode(err error, statusCode int) error {
	return New(err, 0, statusCode)
}

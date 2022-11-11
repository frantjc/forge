package errbubble

import (
	"net/http"

	"github.com/frantjc/go-fn"
)

type Bubble struct {
	error
	exitCode, httpStatusCode int
}

type BubbleOpt func(*Bubble)

func WithHTTPStatusCode(statusCode int) BubbleOpt {
	return func(b *Bubble) {
		b.httpStatusCode = statusCode
	}
}

func WithExitCode(exitCode int) BubbleOpt {
	return func(b *Bubble) {
		b.exitCode = exitCode
	}
}

func (b *Bubble) ExitCode() int {
	if b == nil {
		return 0
	}

	return b.exitCode
}

func (b *Bubble) HTTPStatusCode() int {
	if b == nil {
		return http.StatusOK
	}

	return b.httpStatusCode
}

func New(err error, opts ...BubbleOpt) (b *Bubble) {
	if e, ok := err.(*Bubble); ok {
		e.error = err
		b = e
	} else {
		b = &Bubble{
			err,
			fn.Ternary(err == nil, 0, 1),
			fn.Ternary(err == nil, http.StatusOK, http.StatusInternalServerError),
		}
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

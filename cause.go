package errors

import (
	"fmt"
	"io"
)

// NewCause returns a new error which will have specified cause and
// an origin error which will be used when formatting to string.
func NewCause(cause error, origin error) error {
	return &causeWithOrigin{
		cause,
		origin,
	}
}

type causeWithOrigin struct {
	cause  error
	origin error
}

func (c *causeWithOrigin) Error() string {
	return c.cause.Error() + ": " + c.origin.Error()
}

func (c *causeWithOrigin) Cause() error {
	return c.cause
}

func (c *causeWithOrigin) Origin() error {
	return c.origin
}

func (c *causeWithOrigin) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			switch e := c.cause.(type) {
			case fmt.Formatter:
				e.Format(s, verb)
			}
			io.WriteString(s, "\n")
			switch e := c.origin.(type) {
			case fmt.Formatter:
				e.Format(s, verb)
			}
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, c.Error())
	}
}

// HasCause returns true if one of origins and causes of the error has specified cause.
func HasCause(err error, cause error) bool {
	return HasMatchingCause(err, func(candidate error) bool {
		return candidate == cause
	})
}

// HasMatchingCause returns true if one of origins and causes of the error has matching cause.
func HasMatchingCause(err error, match func(error) bool) bool {
	type originer interface {
		Origin() error
	}

	for err != nil {
		cause := Cause(err)
		if match(cause) {
			return true
		}

		originer, ok := cause.(originer)
		if !ok {
			break
		}

		err = originer.Origin()
	}

	return false
}

package errors

import (
	"fmt"
	"io"
)

// Fields is a convenience alias for map[string]interface{}.
type Fields = map[string]interface{}

// WithField returns a new error with the specified field associated with it.
func WithField(err error, key string, value interface{}) error {
	return WithFields(err, Fields{key: value})
}

// WithFields returns a new error with the specified fields associated with it.
func WithFields(err error, fields Fields) error {
	return &withFields{
		err,
		fields,
	}
}

type withFields struct {
	error
	fields Fields
}

func (w *withFields) Cause() error   { return w.error }
func (w *withFields) Fields() Fields { return w.fields }

func (w *withFields) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			for k, v := range w.fields {
				fmt.Fprintln(s, "{}: {}", k, v)
			}
			e, ok := w.error.(fmt.Formatter)
			if ok {
				e.Format(s, verb)
				return
			}
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

// CollectFields collects all fields until cause error is reached.
func CollectFields(err error) Fields {
	result := Fields{}

	type causer interface {
		Cause() error
	}

	type fielder interface {
		Fields() Fields
	}

	for err != nil {
		fielder, ok := err.(fielder)
		if ok {
			for k, v := range fielder.Fields() {
				if _, ok := result[k]; !ok {
					result[k] = v
				}
			}
		}

		causer, ok := err.(causer)
		if !ok {
			break
		}

		err = causer.Cause()
	}

	return result
}

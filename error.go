package serrs

import (
	"errors"
	"fmt"
)

var _ error = (*simpleError)(nil)
var _ fmt.Formatter = (*simpleError)(nil)
var _ fmt.Stringer = (*simpleError)(nil)
var _ fmt.GoStringer = (*simpleError)(nil)

// simpleError is a simple implementation of the error interface.
type simpleError struct {
	// message is the error message
	message string

	// code is the error code
	// optional
	code Code

	// cause is the cause of the error
	cause error

	// frame is the location where the error occurred
	frame Frame

	// data is the custom data attached to the error
	data CustomData
}

func newSimpleError(msg string, skip int) *simpleError {
	e := new(simpleError)
	e.message = msg
	e.frame = caller(skip + 1)
	return e
}

func asSimpleError(err error) *simpleError {
	if err == nil {
		return nil
	}

	var e *simpleError
	if errors.As(err, &e) {
		return e
	}
	return nil
}

func (s *simpleError) Error() string {
	if s.cause == nil {
		if s.message != "" {
			return s.message
		}
		return ""
	}
	if s.message == "" {
		return s.cause.Error()
	}
	return s.message + ": " + s.cause.Error()
}

func (s *simpleError) Is(target error) bool {
	if targetErr := asSimpleError(target); targetErr != nil {
		return targetErr.getCode() == s.getCode()
	}
	return s == target
}

func (s *simpleError) Unwrap() error {
	return s.cause
}

func (s *simpleError) String() string {
	return s.Error()
}

func (s *simpleError) GoString() string {
	return s.Error()
}

func (s *simpleError) errorSurface() string {
	if s == nil {
		return ""
	}
	if s.message != "" {
		return s.message
	}
	if cause := asSimpleError(s.cause); cause != nil {
		return cause.errorSurface()
	} else if s.cause != nil {
		return s.cause.Error()
	}
	return ""
}

func (s *simpleError) withCode(code Code) error {
	if s == nil {
		return nil
	}
	s.code = code
	return s
}

func (s *simpleError) getCode() Code {
	if s == nil {
		return nil
	}
	if s.code != nil {
		return s.code
	}
	if next := asSimpleError(s.cause); next != nil {
		return next.getCode()
	}
	return nil
}

func (s *simpleError) withData(data CustomData) error {
	if s == nil {
		return nil
	}

	s.data = data
	return s
}

func (s *simpleError) withMessage(msg string) error {
	if s == nil {
		return nil
	}

	s.message = msg
	return s
}

package serrs

import (
	pkgErrors "github.com/pkg/errors"
)

// StackTrace is a method to get the stack trace of the error for sentry-go
func (s *simpleError) StackTrace() pkgErrors.StackTrace {

	f := make([]pkgErrors.Frame, 0, 30)

	if next := asSimpleError(s.cause); next != nil {
		f = append(f, next.StackTrace()...)
	}

	for _, v := range s.frame.frames {
		f = append(f, pkgErrors.Frame(v))
	}

	return f
}

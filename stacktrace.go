package serrs

import (
	"fmt"
)

// StackTrace is a method to get the stack trace of the error for sentry-go
func (s *simpleError) StackTrace() []uintptr {

	f := make([]uintptr, 0, 30)

	if cause := asSimpleError(s.cause); cause != nil {
		f = append(f, cause.StackTrace()...)
	}

	for i, v := range s.frame.frames {
		fmt.Println(i, v)
		f = append(f, v)
	}

	return f
}

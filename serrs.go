package serrs

import (
	"errors"
	"strings"
)

// New creates a new error with the given message and code.
func New(code Code, msg string) error {
	e := newSimpleError(msg, 1).withCode(code)
	return e
}

// Wrap returns an error with a stack trace.
// wrapper is a function that adds information to the error.
func Wrap(err error, ws ...wrapper) error {
	if err == nil {
		return nil
	}
	e := newSimpleError("", 1)
	e.cause = err
	for _, w := range ws {
		w.wrap(e)
	}
	return e
}

// Is reports whether the error is the same as the target error.
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// Origin returns the original error.
func Origin(err error) error {
	for {
		e := Unwrap(err)
		if e == nil {
			return err
		}
		err = e
	}
}

// Unwrap returns the cause of the error.
func Unwrap(e error) error {
	return errors.Unwrap(e)
}

// GetCustomData returns the custom data attached to the error.
func GetCustomData(err error) []CustomData {
	var ss []CustomData
	for {
		if e := asSimpleError(err); e != nil && e.data != nil {
			ss = append(ss, e.data)
		}
		unwrapped := Unwrap(err)
		if unwrapped == nil {
			break
		}
		err = unwrapped
	}
	return ss
}

// GetErrorCode returns the error code attached to the error.
func GetErrorCode(err error) (Code, bool) {
	if e := asSimpleError(err); e != nil {
		if code := e.getCode(); code != nil {
			return code, true
		}
	}
	return nil, false
}

// GetErrorCodeString returns the error code string attached to the error.
func GetErrorCodeString(err error) string {
	code, ok := GetErrorCode(err)
	if !ok {
		return ""
	}
	return code.ErrorCode()
}

// ErrorSurface returns the surface of err.Error()
func ErrorSurface(err error) string {
	if err == nil {
		return ""
	}
	if e := asSimpleError(err); e != nil {
		return e.errorSurface()
	}
	jj := strings.Split(err.Error(), ":")
	if len(jj) == 0 {
		return ""
	}
	return strings.TrimSpace(jj[0])
}

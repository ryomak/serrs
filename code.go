package serrs

// Code is an interface that represents an error code.
type Code interface {
	ErrorCode() string
}

// DefaultCode is a type that represents an error code as a string.
type DefaultCode string

func (s DefaultCode) ErrorCode() string {
	return string(s)
}

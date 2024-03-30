package serrs

// Code is an interface that represents an error code.
type Code interface {
	ErrorCode() string
}

// StringCode is a type that represents an error code as a string.
type StringCode string

func (s StringCode) ErrorCode() string {
	return string(s)
}

const (
	StringCodeUnexpected StringCode = "unexpected"
)

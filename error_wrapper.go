package serrs

// errWrapper is a function that adds information to the error.
type errWrapper interface {
	wrap(err *simpleError)
}

// WithCode returns an error wrapper that adds a code to the error.
func WithCode(code Code) errWrapper {
	return codeWrapper{code: code}
}

type codeWrapper struct {
	code Code
}

func (c codeWrapper) wrap(err *simpleError) {
	_ = err.withCode(c.code)
}

// WithMessage returns an error wrapper that adds a message to the error.
func WithMessage(msg string) errWrapper {
	return messageWrapper{message: msg}
}

type messageWrapper struct {
	message string
}

func (m messageWrapper) wrap(err *simpleError) {
	_ = err.withMessage(m.message)
}

// WithData returns an error wrapper that adds custom data to the error.
func WithData(data CustomData) errWrapper {
	return customDataWrapper{data: data}
}

type customDataWrapper struct {
	data CustomData
}

func (c customDataWrapper) wrap(err *simpleError) {
	_ = err.withData(c.data)
}

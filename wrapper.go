package serrs

type wrapper interface {
	wrap(err *simpleError)
}

func WithCode(code Code) wrapper {
	return codeWrapper{code: code}
}

type codeWrapper struct {
	code Code
}

func (c codeWrapper) wrap(err *simpleError) {
	_ = err.withCode(c.code)
}

func WithMessage(msg string) wrapper {
	return messageWrapper{message: msg}
}

type messageWrapper struct {
	message string
}

func (m messageWrapper) wrap(err *simpleError) {
	_ = err.withMessage(m.message)
}

func WithCustomData(data CustomData) wrapper {
	return customDataWrapper{data: data}
}

type customDataWrapper struct {
	data CustomData
}

func (c customDataWrapper) wrap(err *simpleError) {
	_ = err.withData(c.data)
}

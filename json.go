package serrs

var stackedErrorJsonFormatter = func(msg string, data CustomData) any {
	return defaultErrorJson{
		Message: msg,
		Data:    data,
	}
}

// SetStackedErrorJsonFormatter sets the formatter for StackedErrorJson.
// change the format of the automatically set error
func SetStackedErrorJsonFormatter(f StackedErrorJsonFormatter) {
	stackedErrorJsonFormatter = f
}

// StackedErrorJson returns a slice of JSON objects representing the error stack.
func StackedErrorJson(err error) []any {
	if err == nil {
		return nil
	}
	m := make([]any, 0, 30)

	se := asSimpleError(err)
	if se == nil {
		return append(m, stackedErrorJsonFormatter(err.Error(), nil))
	}
	if causeErr := asSimpleError(se.cause); causeErr != nil {
		m = append(m, StackedErrorJson(causeErr)...)
	} else if se.cause != nil {
		m = append(m, StackedErrorJson(se.cause)...)
	}
	if se.message == "" && se.data == nil {
		return m
	}
	return append(m, stackedErrorJsonFormatter(se.message, se.data))
}

type StackedErrorJsonFormatter func(msg string, data CustomData) any

type defaultErrorJson struct {
	Message string     `json:"message"`
	Data    CustomData `json:"data"`
}

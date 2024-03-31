package serrs

var stackedErrorJsonFormatter = func(err error) any {
	if err == nil {
		return defaultErrorJson{}
	}

	e := asSimpleError(err)
	if e == nil {
		return defaultErrorJson{
			Message: err.Error(),
			Data:    nil,
		}
	}

	return defaultErrorJson{
		Message: e.message,
		Data:    e.data,
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
		return append(m, stackedErrorJsonFormatter(err))
	}
	if causeErr := asSimpleError(se.cause); causeErr != nil {
		m = append(m, StackedErrorJson(causeErr)...)
	} else if se.cause != nil {
		m = append(m, StackedErrorJson(se.cause)...)
	}
	if se.message == "" && se.data == nil {
		return m
	}
	return append(m, stackedErrorJsonFormatter(se))
}

type StackedErrorJsonFormatter func(err error) any

type defaultErrorJson struct {
	Message string     `json:"message"`
	Data    CustomData `json:"data"`
}

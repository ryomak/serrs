package serrs

var stackedErrorJsonFormatter = func(msg string, data CustomData) any {
	return defaultErrorJson{
		Message: msg,
		Data:    data,
	}
}

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
	if ee := asSimpleError(se.cause); ee != nil {
		m = append(m, StackedErrorJson(ee)...)
	} else if se.cause != nil {
		m = append(m, StackedErrorJson(se.cause)...)
	}
	return append(m, stackedErrorJsonFormatter(se.message, se.data))
}

type StackedErrorJsonFormatter func(msg string, data CustomData) any

type defaultErrorJson struct {
	Message string     `json:"message"`
	Data    CustomData `json:"data"`
}

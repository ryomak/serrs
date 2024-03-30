# serrs


[![Go Reference](https://pkg.go.dev/badge/github.com/ryomak/serrs.svg)](https://pkg.go.dev/github.com/ryomak/serrs)
[![GitHub Actions](https://github.com/ryomak/serrs/workflows/test/badge.svg)](https://github.com/ryomak/serrs/actions?query=workflows%3Atest)
[![codecov](https://codecov.io/gh/ryomak/serrs/branch/master/graph/badge.svg)](https://codecov.io/gh/ryomak/serrs)
[![Go Report Card](https://goreportcard.com/badge/github.com/ryomak/serrs)](https://goreportcard.com/report/github.com/ryomak/serrs)


# description

serrs is a package that provides a simple error handling mechanism.



# Installation

```bash
go get -u github.com/ryomak/serrs
```

# Usage
## Create an error
```go

var HogeError = serrs.New(serrs.Unexpceted,"unexpected error")
```

## Wrap an error and add a stack trace
```go

if err := DoSomething(); err != nil {
    // This point is recorded
    return serrs.Wrap(err)
}

fmt.Printf("%+v",err)

// Output:
// - file: ./serrs/format_test.go:22
//   function: github.com/ryomak/serrs_test.TestSerrs_Format
//   msg: wrap error
// - file: ./serrs/format_test.go:14
//   function: github.com/ryomak/serrs_test.TestSerrs_Format
//   code: demo
//   data: {key1:value1,key2:value2}
// - error1
```

## SendSentry
serrs supports sending errors to Sentry.
```go
serrs.ReportSentry(
	err,
	serrs.WithSentryContexts(map[string]sentry.Context{
		"custom": map[string]any{
			"key": "value",
		},
	}),
	serrs.WithSentryTags(map[string]string{
		"code": serrs.GetErrorCodeString(err),
	}),
	serrs.WithSentryLevel(sentry.LevelInfo),
)
```

or 

```go
event := serrs.GenerateSentryEvent(err)
sentry.CaptureEvent(event)
```
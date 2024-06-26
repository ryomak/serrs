# serrs


[![Go Reference](https://pkg.go.dev/badge/github.com/ryomak/serrs.svg)](https://pkg.go.dev/github.com/ryomak/serrs)
[![GitHub Actions](https://github.com/ryomak/serrs/workflows/test/badge.svg)](https://github.com/ryomak/serrs/actions?query=workflows%3Atest)
[![codecov](https://codecov.io/gh/ryomak/serrs/branch/main/graph/badge.svg)](https://codecov.io/gh/ryomak/serrs)
[![Go Report Card](https://goreportcard.com/badge/github.com/ryomak/serrs)](https://goreportcard.com/report/github.com/ryomak/serrs)


# description

serrs is a library designed to simplify error handling in your applications.
By using serrs, developers can effortlessly manage stack traces and integrate with monitoring tools like Sentry. 

# Installation

```bash
go get -u github.com/ryomak/serrs
```

# Usage
## Create an error
```go

var HogeError = serrs.Wrap(err)
```

or 

```go
var InvalidParameterError = serrs.New(serrs.DefaultCode("invalid_parameter"),"invalid parameter error")
```

## Wrap an error and add a stack trace
```go

if err := DoSomething(); err != nil {
    // This point is recorded
    return serrs.Wrap(err)
}

fmt.Printf("%+v",err)

// Output Example:
// - file: ./serrs/format_test.go:22
//   function: github.com/ryomak/serrs_test.TestSerrs_Format
//   msg: 
// - file: ./serrs/format_test.go:14
//   function: github.com/ryomak/serrs_test.TestSerrs_Format
//   code: demo
//   data: {key1:value1,key2:value2}
// - error1
```

### Parameters
The parameters that can be managed with serrs are three: `code`, `message`, and `data`.
`WithXXX` functions can be used to add additional information to the error.
- code: error code
- message: err text
- data: custom data

**data**  
`data` is a custom data that can be added to the error. The data is output to the log.  
If the type satisfies the CustomData interface, any type can be added.

```go
if err := DoSomething(); err != nil {
    return serrs.Wrap(err, serrs.WithData(serrs.DefaultCustomData{
        "key": "value",
    }))
}
```

### Get Additional Data Functions
- GetCustomData(err error): Get custom data from error
- GetErrorCode(err error): Get error code from error
- GetErrorCodeString(err error): Get error code string from error
- ErrorSurface(err error): Get top level error message
- Origin(err error): Get original error

## check error match
```go
var HogeError = serrs.New(serrs.DefaultCode("unexpected"),"unexpected error")

if serrs.Is(HogeError) {
    
}
```

## Send Sentry
supports sending reports to Sentry.
The location where serrs.Wrap is executed is saved as a stack trace and displayed cleanly on Sentry. In addition, any added custom data or messages are also displayed as additional data on Sentry.

```go

serrs.ReportSentry(
    err, 
    // Customize the contexts 
    serrs.WithSentryContexts(map[string]sentry.Context{
        "custom": map[string]any{
            "key": "value",
        },
    }), 
    // Customize the Sentry tags 
    serrs.WithSentryTags(map[string]string{
        "code": serrs.GetErrorCodeString(err),
    }), 
    // Customize the Sentry Level 
    serrs.WithSentryLevel(sentry.LevelInfo),
)
```

or 

```go

import (
    "github.com/getsentry/sentry-go"
)

func main() {
    sentry.Init(sentry.ClientOptions{
        Dsn: "your-dsn",
    })
    defer sentry.Flush(2 * time.Second)
	
    if err := DoSomething(); err != nil {
        event := serrs.GenerateSentryEvent(err)
        sentry.CaptureEvent(event)
    }
}
```

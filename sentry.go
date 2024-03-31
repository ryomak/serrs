package serrs

import (
	"context"
	"fmt"

	sentry "github.com/getsentry/sentry-go"
)

// StackTrace is a method to get the stack trace of the error for sentry-go
// https://github.com/getsentry/sentry-go/blob/master/stacktrace.go#L84-L87
func (s *simpleError) StackTrace() []uintptr {
	origin := originSimpleError(s)
	if origin == nil {
		return []uintptr{}
	}
	if len(origin.frames) <= 1 {
		return origin.frames
	}
	return origin.frames[1:]
}

// GenerateSentryEvent is a method to generate a sentry event from an error
func GenerateSentryEvent(err error, ws ...sentryWrapper) *sentry.Event {
	if err == nil {
		return nil
	}
	exceptionType := Origin(err).Error()
	errCode, ok := GetErrorCode(err)
	if !ok {
		errCode = DefaultCode("unknown")
	} else {
		exceptionType = fmt.Sprintf("[%s] %s", errCode, exceptionType)
	}

	event := sentry.NewEvent()
	event.Level = sentry.LevelError
	event.Exception = []sentry.Exception{{
		Value:      err.Error(),
		Type:       exceptionType,
		Stacktrace: sentry.ExtractStacktrace(err),
	}}
	event.Contexts = map[string]sentry.Context{
		"error detail": {
			"history": StackedErrorJson(err),
			"code":    errCode.ErrorCode(),
		},
	}

	for _, w := range ws {
		w.wrap(event)
	}

	return event
}

// ReportSentry is a method to report an error to sentry
func ReportSentry(err error, ws ...sentryWrapper) {
	event := GenerateSentryEvent(err, ws...)
	sentry.CaptureEvent(event)
}

// ReportSentryWithContext is a method to report an error to sentry with a context
func ReportSentryWithContext(ctx context.Context, err error, ws ...sentryWrapper) {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		ReportSentry(err, ws...)
		return
	}
	event := GenerateSentryEvent(err, ws...)
	hub.CaptureEvent(event)
}

type sentryWrapper interface {
	wrap(event *sentry.Event) *sentry.Event
}

// WithSentryContexts is a function to add contexts to a sentry event
func WithSentryContexts(m map[string]sentry.Context) sentryWrapper {
	return sentryEventContextWrapper{m}
}

type sentryEventContextWrapper struct {
	m map[string]sentry.Context
}

func (s sentryEventContextWrapper) wrap(event *sentry.Event) *sentry.Event {
	for k, v := range s.m {
		event.Contexts[k] = v
	}

	return event
}

// WithSentryTags is a function to add tags to a sentry event
func WithSentryTags(m map[string]string) sentryWrapper {
	return sentryEventTagWrapper{m}
}

type sentryEventTagWrapper struct {
	m map[string]string
}

func (s sentryEventTagWrapper) wrap(event *sentry.Event) *sentry.Event {
	for k, v := range s.m {
		event.Tags[k] = v
	}

	return event
}

// WithSentryLevel is a function to set the level of a sentry event
func WithSentryLevel(l sentry.Level) sentryWrapper {
	return sentryEventLevelWrapper{l}
}

type sentryEventLevelWrapper struct {
	l sentry.Level
}

func (s sentryEventLevelWrapper) wrap(event *sentry.Event) *sentry.Event {
	event.Level = s.l

	return event
}

func originSimpleError(err error) *simpleError {
	var e *simpleError
	for {
		if err == nil {
			return e
		}
		if ee := asSimpleError(err); ee != nil {
			e = ee
		}
		err = Unwrap(err)
	}
}

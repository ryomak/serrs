package serrs

import (
	sentry "github.com/getsentry/sentry-go"
)

// StackTrace is a method to get the stack trace of the error for sentry-go
// https://github.com/getsentry/sentry-go/blob/master/stacktrace.go#L84-L87
func (s *simpleError) StackTrace() []uintptr {

	frames := make([]uintptr, 0, 30)
	origin := originSimpleError(s)
	if origin == nil {
		return frames
	}
	for _, frame := range origin.frame.frames {
		frames = append(frames, frame)
	}
	if len(frames) > 0 {
		frames = frames[1:]
	}

	return frames
}

// GenerateSentryEvent is a method to generate a sentry event from an error
func GenerateSentryEvent(err error, ws ...sentryWrapper) *sentry.Event {
	if err == nil {
		return nil
	}
	errCode, ok := GetErrorCode(err)
	if !ok {
		errCode = StringCode("unknown")
	}
	event := sentry.NewEvent()
	event.Level = sentry.LevelError
	event.Exception = []sentry.Exception{{
		Value:      err.Error(),
		Type:       Origin(err).Error(),
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

package serrs_test

import (
	"context"
	"errors"
	"testing"

	sentry "github.com/getsentry/sentry-go"
	"github.com/ryomak/serrs"
)

func TestGenerateSentryEvent_WithNilError(t *testing.T) {
	t.Parallel()

	event := serrs.GenerateSentryEvent(nil)
	if event != nil {
		t.Errorf("Expected nil, but got %v", event)
	}
}

func TestGenerateSentryEvent_WithUnknownErrorCode(t *testing.T) {
	t.Parallel()

	err := errors.New("test error")
	event := serrs.GenerateSentryEvent(err)
	if event.Contexts["error detail"]["code"] != "unknown" {
		t.Errorf("Expected 'unknown', but got %v", event.Contexts["error detail"]["code"])
	}
}

func TestGenerateSentryEvent_WithKnownErrorCode(t *testing.T) {
	t.Parallel()

	err := serrs.Wrap(serrs.New(serrs.DefaultCode("known"), "test error"))
	event := serrs.GenerateSentryEvent(err)
	if event.Contexts["error detail"]["code"] != "known" {
		t.Errorf("Expected 'known', but got %v", event.Contexts["error detail"]["code"])
	}
}

func TestGenerateSentryEvent_WithSentryOptions(t *testing.T) {
	t.Parallel()

	err := serrs.Wrap(serrs.New(serrs.DefaultCode("known"), "test error"))
	event := serrs.GenerateSentryEvent(
		err,
		serrs.WithSentryTags(map[string]string{"key": "value"}),
		serrs.WithSentryLevel(sentry.LevelWarning),
		serrs.WithSentryContexts(map[string]sentry.Context{
			"requestData": map[string]any{
				"key": "value",
			},
		}),
	)
	checkEqual(t, event.Contexts["requestData"], map[string]any{"key": "value"})
	checkEqual(t, event.Level, sentry.LevelWarning)
	checkEqual(t, event.Tags, map[string]string{"key": "value"})
}

func TestReportSentry(t *testing.T) {
	t.Parallel()

	err := serrs.New(serrs.DefaultCode("test"), "test error")
	serrs.ReportSentry(err)
	serrs.ReportSentryWithContext(context.Background(), err)
}

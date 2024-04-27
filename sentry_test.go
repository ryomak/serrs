package serrs_test

import (
	"errors"
	"testing"

	"github.com/ryomak/serrs"
)

func TestGenerateSentryEvent_WithNilError(t *testing.T) {
	event := serrs.GenerateSentryEvent(nil)
	if event != nil {
		t.Errorf("Expected nil, but got %v", event)
	}
}

func TestGenerateSentryEvent_WithUnknownErrorCode(t *testing.T) {
	err := errors.New("test error")
	event := serrs.GenerateSentryEvent(err)
	if event.Contexts["error detail"]["code"] != "unknown" {
		t.Errorf("Expected 'unknown', but got %v", event.Contexts["error detail"]["code"])
	}
}

func TestGenerateSentryEvent_WithKnownErrorCode(t *testing.T) {
	err := serrs.Wrap(serrs.New(serrs.StringCode("known"), "test error"))
	event := serrs.GenerateSentryEvent(err)
	if event.Contexts["error detail"]["code"] != "known" {
		t.Errorf("Expected 'known', but got %v", event.Contexts["error detail"]["code"])
	}
}

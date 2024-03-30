package serrs_test

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/ryomak/serrs"
)

func checkMatch(t *testing.T, target, regexpStr string) {
	t.Helper()
	r := regexp.MustCompile(regexpStr)
	if !r.MatchString(target) {
		t.Errorf("%q does not match %q", target, regexpStr)
	}
}

func checkEqual(t *testing.T, a, b any) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%T(%#v) does not equal to %T(%#v)", a, a, b, b)
	}
}

func checkIsError(t *testing.T, err error, target error) {
	t.Helper()
	if !serrs.Is(err, target) {
		t.Errorf("%v is not %v", err, target)
	}
}

func mustErrorCode(t *testing.T, err error) serrs.Code {
	code, ok := serrs.GetErrorCode(err)
	if !ok {
		t.Errorf("failed to get code from error")
		return nil
	}
	return code
}

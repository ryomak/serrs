package serrs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ryomak/serrs"
)

func TestSerrs_Format(t *testing.T) {

	e1 := errors.New("error1")
	err := serrs.Wrap(
		e1,
		serrs.WithCode(serrs.DefaultCode("demo")),
		serrs.WithData(serrs.DefaultCustomData{
			"key1": "value1",
			"key2": "value2",
		}),
	)
	err = serrs.Wrap(err, serrs.WithMessage("wrap error"))

	checkMatch(t, fmt.Sprintf("%+v", err), `- file: .*serrs\/format_test.go:22
  function: .*serrs_test.TestSerrs_Format
  msg: wrap error
- file: .*format_test.go:14
  function: .*serrs_test.TestSerrs_Format
  code: demo
  data: {key1:value1,key2:value2}
- error1`)

	checkEqual(t, fmt.Sprintf("%v", err), "wrap error: error1")
	checkEqual(t, fmt.Sprintf("%#v", err), "wrap error: error1")
	checkEqual(t, fmt.Sprintf("%s", err), "wrap error: error1")

	// unexpected format
	//nolint
	checkEqual(t, fmt.Sprintf("%a", err), "%!a(*serrs.simpleError)")
}

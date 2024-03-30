package serrs_test

import (
	"testing"

	"github.com/ryomak/serrs"
)

func TestDefaultCustomData_String(t *testing.T) {
	var user struct {
		name string
	}
	user.name = "hogehoge"

	in := serrs.DefaultCustomData{
		"key1": "value1",
		"user": user,
	}
	out := in.String()
	checkEqual(t, out, "{key1:value1,user:{name:hogehoge}}")
}

func TestDefaultCustomData_Clone(t *testing.T) {
	var user struct {
		name string
	}
	value3 := 3
	in := serrs.DefaultCustomData{
		"key1": "value1",
		"user": user,
		"ptr":  &value3,
	}
	out := in.Clone()

	checkEqual(t, in, out)

	delete(in, "key1")
	if len(in) == len(out.(serrs.DefaultCustomData)) {
		t.Errorf("Clone() should return a new map")
	}
}

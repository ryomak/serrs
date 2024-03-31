package serrs_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ryomak/serrs"
)

func TestStackedErrorJson(t *testing.T) {
	tests := []struct {
		name string
		in   error
		want string
	}{
		{
			name: "simpleError",
			in: serrs.Wrap(
				serrs.New(serrs.DefaultCode("unexpected"), "unexpected error"),
				serrs.WithMessage("wrap error"),
				serrs.WithData(serrs.DefaultCustomData{"key1": "value1"}),
			),
			want: `[{"message":"unexpected error","data":null},{"message":"wrap error","data":{"key1":"value1"}}]`,
		},
		{
			name: "pure error",
			in:   errors.New("unexpected error"),
			want: `[{"message":"unexpected error","data":null}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := serrs.StackedErrorJson(tt.in)
			outJson, err := json.Marshal(out)
			if err != nil {
				t.Errorf("json.Marshal() = %v, want %v", outJson, tt.want)
			}
			checkEqual(t, string(outJson), tt.want)
		})
	}
}

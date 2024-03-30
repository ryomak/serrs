package serrs_test

import (
	"testing"

	"github.com/ryomak/serrs"
)

func TestStringCode_ErrorCode(t *testing.T) {
	tests := []struct {
		name string
		s    serrs.StringCode
		want string
	}{
		{
			name: "unexpected",
			s:    serrs.StringCodeUnexpected,
			want: "unexpected",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ErrorCode(); got != tt.want {
				t.Errorf("StringCode.GetErrorCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

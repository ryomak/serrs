package serrs_test

import (
	"testing"

	"github.com/ryomak/serrs"
)

func TestDefaultCode_ErrorCode(t *testing.T) {
	tests := []struct {
		name string
		s    serrs.DefaultCode
		want string
	}{
		{
			name: "unexpected",
			s:    serrs.DefaultCode("unexpected"),
			want: "unexpected",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ErrorCode(); got != tt.want {
				t.Errorf("DefaultCode.GetErrorCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

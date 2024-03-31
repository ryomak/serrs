package serrs_test

import (
	"errors"
	"testing"

	"github.com/ryomak/serrs"
)

func TestSerrs(t *testing.T) {

	baseErr := serrs.New(serrs.DefaultCode("hoge_error"), "hoge error")
	type want struct {
		checkNil bool
		code     serrs.Code
		error    string
	}
	tests := map[string]struct {
		in   error
		want want
	}{
		"New": {
			in: baseErr,
			want: want{
				code:  serrs.DefaultCode("hoge_error"),
				error: "hoge error",
			},
		},
		"Wrap": {
			in: serrs.Wrap(baseErr, serrs.WithMessage("wrap error")),
			want: want{
				code:  serrs.DefaultCode("hoge_error"),
				error: "wrap error: hoge error",
			},
		},
		"Wrap_nil": {
			in: serrs.Wrap(nil, serrs.WithMessage("wrap error")),
			want: want{
				checkNil: true,
				code:     nil,
				error:    "",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.want.checkNil {
				if tt.in != nil {
					t.Errorf("got %v, want nil", tt.in)
				}
				return
			}
			checkEqual(t, mustErrorCode(t, tt.in), tt.want.code)
			checkEqual(t, tt.in.Error(), tt.want.error)
		})
	}
}

func TestSimpleError_Is(t *testing.T) {
	tests := map[string]struct {
		err    error
		target error
		want   bool
	}{
		"simpleError -> simpleError: same code": {
			err:    serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"),
			target: serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"),
			want:   true,
		},
		"simpleError -> simpleError: other code": {
			err:    serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"),
			target: serrs.New(serrs.DefaultCode("fuga_error"), "fuga error"),
			want:   false,
		},
		"wrap simpleError -> simpleError: same code": {
			err:    serrs.Wrap(serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"), serrs.WithMessage("wrap error")),
			target: serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"),
			want:   true,
		},
		"wrap simpleError -> simpleError : same code WithCode": {
			err:    serrs.Wrap(serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"), serrs.WithCode(serrs.DefaultCode("fuga_error"))),
			target: serrs.New(serrs.DefaultCode("fuga_error"), "fuga error"),
			want:   true,
		},
		"normal error -> normal error": {
			err:    errors.ErrUnsupported,
			target: errors.ErrUnsupported,
			want:   true,
		},
		"wrap normal error -> normal error": {
			err:    serrs.Wrap(errors.ErrUnsupported, serrs.WithMessage("wrap error")),
			target: errors.ErrUnsupported,
			want:   true,
		},
		"normal error -> wrap normal error": {
			err:    errors.ErrUnsupported,
			target: serrs.Wrap(errors.ErrUnsupported, serrs.WithMessage("wrap error")),
			want:   false,
		},
		"wrap normal error -> other wrap normal error": {
			err:    serrs.Wrap(errors.ErrUnsupported, serrs.WithCode(serrs.DefaultCode("fuga_error"))),
			target: serrs.Wrap(errors.ErrUnsupported, serrs.WithCode(serrs.DefaultCode("hoge_error"))),
			want:   false,
		},
		"wrap simpleError-> simpleError: same code": {
			err:    serrs.Wrap(serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"), serrs.WithCode(serrs.DefaultCode("fuga_error"))),
			target: serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"),
			want:   true,
		},
		"simpleError -> wrap simpleError: same code": {
			err:    serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"),
			target: serrs.Wrap(serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"), serrs.WithCode(serrs.DefaultCode("fuga_error"))),
			want:   false,
		},
		"wrap normal error -> nil": {
			err:    serrs.Wrap(errors.ErrUnsupported),
			target: nil,
			want:   false,
		},
		"nil -> wrap normal error": {
			err:    nil,
			target: serrs.Wrap(errors.ErrUnsupported),
			want:   false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			checkEqual(t, serrs.Is(tt.err, tt.target), tt.want)
		})
	}
}

func TestOrigin(t *testing.T) {
	tests := map[string]struct {
		in   error
		want error
	}{
		"simple error": {
			in:   serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"),
			want: serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"),
		},
		"wrap error": {
			in:   serrs.Wrap(serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"), serrs.WithMessage("wrap error")),
			want: serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"),
		},
		"nil error": {
			in:   nil,
			want: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			checkIsError(t, serrs.Origin(tt.in), tt.want)
		})
	}
}

func TestWrap_WithCustomData(t *testing.T) {
	tests := map[string]struct {
		in   error
		data serrs.CustomData
		want []serrs.CustomData
	}{
		"simple error": {
			in: serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"),
			data: serrs.DefaultCustomData{
				"key": "value",
			},
			want: []serrs.CustomData{
				serrs.DefaultCustomData{
					"key": "value",
				},
			},
		},
		"nil error": {
			in: nil,
			data: serrs.DefaultCustomData{
				"key": "value",
			},
			want: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			checkEqual(t, serrs.GetCustomData(serrs.Wrap(tt.in, serrs.WithData(tt.data))), tt.want)
		})
	}
}

func TestGetCustomData_Example(t *testing.T) {
	var in error
	if err := func() error {
		return serrs.Wrap(serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"), serrs.WithData(serrs.DefaultCustomData{
			"key": "value",
		}))
	}(); err != nil {
		in = serrs.Wrap(
			err,
			serrs.WithData(serrs.DefaultCustomData{
				"key2": "value2",
			}),
		)
	}

	checkEqual(t, serrs.GetCustomData(in), []serrs.CustomData{
		serrs.DefaultCustomData{
			"key2": "value2",
		},
		serrs.DefaultCustomData{
			"key": "value",
		},
	})
}

func TestErrorSurface(t *testing.T) {
	tests := map[string]struct {
		in   error
		want string
	}{
		"simple error": {
			in:   serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"),
			want: "hoge error",
		},
		"wrap error": {
			in:   serrs.Wrap(serrs.New(serrs.DefaultCode("hoge_error"), "hoge error"), serrs.WithMessage("wrap error")),
			want: "wrap error",
		},
		"wrap error without msg": {
			in:   serrs.Wrap(serrs.New(serrs.DefaultCode("hoge_error"), "hoge error")),
			want: "hoge error",
		},
		"nil error": {
			in:   nil,
			want: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			checkEqual(t, serrs.ErrorSurface(tt.in), tt.want)
		})
	}
}

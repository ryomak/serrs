package serrs

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Format implements fmt.Formatter interface
func (s *simpleError) Format(state fmt.State, v rune) {
	formatError(s, state, v)
}

func (s *simpleError) printerFormat(p printer) error {
	var message string
	if s.code != nil {
		message += fmt.Sprintf("code: %s\n", s.code.ErrorCode())
	}
	if s.message != "" {
		message += fmt.Sprintf("msg: %s\n", s.message)
	}
	if s.data != nil {
		message += fmt.Sprintf("data: %v", s.data)
	}

	// print stack frame
	s.frame.format(p)

	// print message
	p.Print(message)

	// print cause recursively
	return s.cause
}

// formatError formats the error according to the verb and flags.
func formatError(f *simpleError, s fmt.State, verb rune) {
	var (
		sep    = " " // separator before next error
		p      = &state{State: s}
		direct = true
	)

	var err error = f

	switch verb {

	case 'v', 's':
		if s.Flag('#') {
			if stringer, ok := err.(fmt.GoStringer); ok {
				_, _ = io.WriteString(&p.buf, stringer.GoString())
				goto exit
			}
		} else if s.Flag('+') {
			p.printDetail = true
			sep = "- "
		} else {
			if stringer, ok := err.(fmt.Stringer); ok {
				_, _ = io.WriteString(&p.buf, stringer.String())
				goto exit
			}
		}

	case 'q', 'x', 'X':
		direct = false

	default:
		p.buf.WriteString("%!")
		p.buf.WriteRune(verb)
		p.buf.WriteByte('(')
		switch {
		case err != nil:
			p.buf.WriteString(reflect.TypeOf(f).String())
		default:
			p.buf.WriteString("<nil>")
		}
		p.buf.WriteByte(')')
		_, _ = io.Copy(s, &p.buf)
		return
	}

loop:
	for {
		p.buf.WriteString(sep)
		switch v := err.(type) {
		case *simpleError:
			err = v.printerFormat(p)
		case fmt.Formatter:
			v.Format(p, 'v')
			break loop
		default:
			_, _ = io.WriteString(&p.buf, v.Error())
			break loop
		}
		if err == nil {
			break loop
		}
		if p.printDetail {
			p.buf.WriteString("\n")
		}
		if p.needColon || !p.printDetail {
			p.buf.WriteByte(':')
			p.needColon = false
		}
		p.inDetail = false
		p.needNewline = false
	}

exit:
	width, okW := s.Width()
	prec, okP := s.Precision()

	if !direct || (okW && width > 0) || okP {
		// Construct format string from State s.
		format := []byte{'%'}
		if s.Flag('-') {
			format = append(format, '-')
		}
		if s.Flag('+') {
			format = append(format, '+')
		}
		if s.Flag(' ') {
			format = append(format, ' ')
		}
		if okW {
			format = strconv.AppendInt(format, int64(width), 10)
		}
		if okP {
			format = append(format, '.')
			format = strconv.AppendInt(format, int64(prec), 10)
		}
		format = append(format, string(verb)...)
		_, _ = fmt.Fprintf(s, string(format), strings.TrimSpace(p.buf.String()))
	} else {
		_, _ = io.Copy(s, &p.buf)
	}
}

var detailSep = []byte("\n  ")

var _ printer = (*state)(nil)

// state tracks error printing state. It implements fmt.State.
type state struct {
	fmt.State
	buf bytes.Buffer

	printDetail bool
	inDetail    bool
	needColon   bool
	needNewline bool
}

func (s *state) Write(b []byte) (n int, err error) {
	if s.printDetail {
		if len(b) == 0 {
			return 0, nil
		}
		if s.inDetail && s.needColon {
			s.needNewline = true
			if b[0] == '\n' {
				b = b[1:]
			}
		}
		k := 0
		for i, c := range b {
			if s.needNewline {
				if s.inDetail && s.needColon {
					s.buf.WriteByte(':')
					s.needColon = false
				}
				s.buf.Write(detailSep)
				s.needNewline = false
			}
			if c == '\n' {
				s.buf.Write(b[k:i])
				k = i + 1
				s.needNewline = true
			}
		}
		s.buf.Write(b[k:])
		if !s.inDetail {
			s.needColon = true
		}
	} else if !s.inDetail {
		s.buf.Write(b)
	}
	return len(b), nil
}

func (s *state) Print(args ...interface{}) {
	if !s.inDetail || s.printDetail {
		_, _ = fmt.Fprint(s, args...)
	}
}

func (s *state) Printf(format string, args ...interface{}) {
	if !s.inDetail || s.printDetail {
		_, _ = fmt.Fprintf(s, format, args...)
	}
}

func (s *state) Detail() bool {
	s.inDetail = true
	return s.printDetail
}

// printer is the interface that wraps the basic Print and Printf methods.
type printer interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Detail() bool
}

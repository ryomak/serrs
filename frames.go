package serrs

import (
	"runtime"
)

// A Frames represents a program counter inside a stack frames.
type Frames []uintptr

// caller returns a Frames that describes a frames on the caller's stack.
func caller(skip int) Frames {
	f := [32]uintptr{}
	n := runtime.Callers(skip+1, f[:])
	return f[:n]
}

// location returns the function, file, and line number of a Frames.
func (f Frames) location() (function, file string, line int) {
	frames := runtime.CallersFrames(f[:])
	if _, ok := frames.Next(); !ok {
		return "", "", 0
	}
	fr, ok := frames.Next()
	if !ok {
		return "", "", 0
	}
	return fr.Function, fr.File, fr.Line
}

func (f Frames) format(p printer) {
	if p.Detail() {
		function, file, line := f.location()
		if file != "" {
			p.Printf("file: %s:%d\n", file, line)
		}
		if function != "" {
			p.Printf("function: %s\n", function)
		}
	}
}

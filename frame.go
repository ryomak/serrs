package serrs

import (
	"runtime"
)

// A Frame represents a program counter inside a stack frame.
type Frame struct {
	// https://go.googlesource.com/go/+/032678e0fb/src/runtime/extern.go#169
	frames [3]uintptr
}

// caller returns a Frame that describes a frame on the caller's stack.
func caller(skip int) Frame {
	var s Frame
	runtime.Callers(skip+1, s.frames[:])
	return s
}

// location returns the function, file, and line number of a Frame.
func (f Frame) location() (function, file string, line int) {
	frames := runtime.CallersFrames(f.frames[:])
	if _, ok := frames.Next(); !ok {
		return "", "", 0
	}
	fr, ok := frames.Next()
	if !ok {
		return "", "", 0
	}
	return fr.Function, fr.File, fr.Line
}

func (f Frame) format(p printer) {
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

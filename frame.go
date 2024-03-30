package serrs

import (
	"runtime"
)

// A Frame スタックフレームの情報を保持する構造体
type Frame struct {
	// https://go.googlesource.com/go/+/032678e0fb/src/runtime/extern.go#169
	frames [3]uintptr
}

// caller 呼び出し元情報のデータを返す。(メモリのアドレス/ファイル名/行番号など)
// skipは呼び出し元の情報をスキップするスタックフレームの数。0は呼び出し元, 1は呼び出し元の呼び出し元が返却される
func caller(skip int) Frame {
	var s Frame
	runtime.Callers(skip+1, s.frames[:])
	return s
}

// location frameから関数、ファイル、行を返す
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

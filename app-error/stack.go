package apperror

import (
	"runtime"
)

const (
	maximumStackDeep = 32
)

// stack represents a stack of program counters.
type stack []uintptr

func (s *stack) Frames() *runtime.Frames {
	return runtime.CallersFrames(*s)
}

func callers(skip int) *stack {
	var pcs [maximumStackDeep]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

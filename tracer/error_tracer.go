package tracer

import (
	"fmt"
	"runtime"
)

const (
	stackDepthLength = 20
	skippedCals      = 3 // call for skip top frames
)

type ErrorTracer interface {
	Trace() *runtime.Frames
	Unwrap() error
	Error() string
}

type tError struct {
	msg string

	stack   *runtime.Frames
	wrapped error
}

func (e *tError) Error() string {
	return e.msg
}

func (e *tError) Trace() *runtime.Frames {
	return e.stack
}

func (e *tError) Unwrap() error {
	return e.wrapped
}

// Error create traced error with message from other error
func Error(msg string, err error) error {
	return &tError{
		msg:     msg,
		wrapped: err,
		stack:   CatchStack(skippedCals),
	}
}

// New create new traced error with provided message
func New(msg string) error {
	return &tError{
		msg:   msg,
		stack: CatchStack(skippedCals),
	}
}

// Errorf create traced error with formatted error message
func Errorf(format string, args ...any) error {
	formatted := fmt.Errorf(format, args...)
	return &tError{
		msg:     formatted.Error(),
		wrapped: formatted,
		stack:   CatchStack(skippedCals),
	}
}

// Wraps wrap provided error to traced error
func Wraps(err error) error {
	return &tError{
		msg:     err.Error(),
		wrapped: err,
		stack:   CatchStack(skippedCals),
	}
}

// CatchStack Catch trace stack for tracing
func CatchStack(skip int) *runtime.Frames {
	rpc := make([]uintptr, stackDepthLength)
	n := runtime.Callers(skip, rpc)
	if n < 1 {
		return nil
	}
	return runtime.CallersFrames(rpc)
}

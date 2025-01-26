package tracer

import (
	"errors"
	"runtime"
	"testing"
)

func unwrapFrames(frames *runtime.Frames) []runtime.Frame {
	frame, more := frames.Next()
	res := []runtime.Frame{frame}
	for more {
		frame, more = frames.Next()
		res = append(res, frame)
	}
	return res
}

func assertTraces(t *testing.T, expected, actual []runtime.Frame) {
	if len(expected) != len(actual) {
		t.Fatalf("expected %d frames, got %d", len(expected), len(actual))
	}

	for i, e := range expected {
		a := actual[i]
		if e.Function != a.Function {
			t.Fatalf("traces not match %v, %v actual", expected, actual)
		}
	}
}

func TestError_New(t *testing.T) {
	msg := "error message"
	err := New(msg)

	if err.Error() != msg {
		t.Fatalf("Error message does not match: expected %s, got %s", msg, err.Error())
	}

	tErr, ok := err.(ErrorTracer)
	if !ok {
		t.Fatalf("Error is not TracedError")
	}
	expectedTrace := unwrapFrames(CatchStack(2))
	actualTrace := unwrapFrames(tErr.Trace())

	assertTraces(t, expectedTrace, actualTrace)
}

func TestError_Error(t *testing.T) {
	msg := "error message"
	parrentErr := errors.New("some error")
	err := Error(msg, parrentErr)

	if err.Error() != msg {
		t.Fatalf("Error message does not match: expected %s, got %s", msg, err.Error())
	}

	if errors.Unwrap(err) != parrentErr {
		t.Fatalf("Error message does not match: expected %v, got %v", parrentErr.Error(), err.Error())
	}

	tErr, ok := err.(ErrorTracer)
	if !ok {
		t.Fatalf("Error is not TracedError")
	}
	expectedTrace := unwrapFrames(CatchStack(2))
	actualTrace := unwrapFrames(tErr.Trace())

	assertTraces(t, expectedTrace, actualTrace)
}

func TestError_Errorf(t *testing.T) {
	msg := "error message %w"
	formattedMsg := "error message some error"
	parrentErr := errors.New("some error")
	err := Errorf(msg, parrentErr)

	if err.Error() != formattedMsg {
		t.Fatalf("Error message does not match: expected %s, got %s", formattedMsg, err.Error())
	}

	if !errors.Is(err, parrentErr) {
		t.Fatalf("Error message does not match: expected %v, got %v", parrentErr.Error(), err.Error())
	}

	tErr, ok := err.(ErrorTracer)
	if !ok {
		t.Fatalf("Error is not TracedError")
	}
	expectedTrace := unwrapFrames(CatchStack(2))
	actualTrace := unwrapFrames(tErr.Trace())

	assertTraces(t, expectedTrace, actualTrace)
}

func TestError_Wraps(t *testing.T) {
	parrentErr := errors.New("some error")
	err := Wraps(parrentErr)

	if err.Error() != parrentErr.Error() {
		t.Fatalf("Error message does not match: expected %s, got %s", parrentErr.Error(), err.Error())
	}

	if !errors.Is(err, parrentErr) {
		t.Fatalf("Error message does not match: expected %v, got %v", parrentErr.Error(), err.Error())
	}

	tErr, ok := err.(ErrorTracer)
	if !ok {
		t.Fatalf("Error is not TracedError")
	}
	expectedTrace := unwrapFrames(CatchStack(2))
	actualTrace := unwrapFrames(tErr.Trace())

	assertTraces(t, expectedTrace, actualTrace)
}

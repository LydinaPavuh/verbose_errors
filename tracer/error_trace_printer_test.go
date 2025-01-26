package tracer

import (
	"fmt"
	"strings"
	"testing"
)

func A() error {
	return New("Some error")
}

func B() error {
	return A()
}

func TestPrintTrace(t *testing.T) {
	err := B()
	res := PrintTrace(err)

	expectedLines := []string{
		"ERROR: Some error\n",
		"tracer.B",
		"tracer.A",
		"tracer.TestPrintTrace",
	}

	for _, line := range expectedLines {
		if !strings.Contains(res, line) {
			t.Fatalf("Expected to find %q in %s", line, res)
		}
	}

}

func TestPrintTrace_WithCallers(t *testing.T) {
	err := B()
	res := PrintTraceWithOpts(err, &FormatOpts{
		WithTrace:         false,
		WithCaller:        true,
		WithUntracedWraps: false,
	})

	expectedLines := []string{
		"ERROR: Some error\n",
		"tracer.A",
	}

	for _, line := range expectedLines {
		if !strings.Contains(res, line) {
			t.Fatalf("Expected to find %q in %s", line, res)
		}
	}
}

func TestPrintTrace_WithUntracedWraps(t *testing.T) {
	err := New("root error")
	err = fmt.Errorf("some untraced error %w", err)
	err = Error("Some error", err)

	res := PrintTraceWithOpts(err, &FormatOpts{
		WithTrace:         false,
		WithCaller:        false,
		WithUntracedWraps: true,
	})

	expectedLines := []string{
		"ERROR: Some error",
		"ERROR: root error",
		"UNTRACED WRAPS:",
		"0. some untraced error root error",
	}

	for _, line := range expectedLines {
		if !strings.Contains(res, line) {
			t.Fatalf("Expected to find %q in %s", line, res)
		}
	}
}

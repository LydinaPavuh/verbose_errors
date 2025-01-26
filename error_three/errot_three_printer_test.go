package error_three

import (
	"errors"
	"testing"
)

func TestPrintErrorThree(t *testing.T) {
	root := errors.New("root error")
	second := Error("second error", root)

	res := PrintErrorThree(second)
	expected := "-- 1: second error\n---- 2: root error\n"
	if res != expected {
		t.Fatalf("PrintErrorThree(%#v) returned %#v, expected %#v", second, res, expected)
	}

	third := Error("third error", second, errors.New("additional"))
	expected = "-- 1: third error\n---- 2: additional\n---- 2: second error\n------ 3: root error\n"
	res = PrintErrorThree(third)
	if res != expected {
		t.Fatalf("PrintErrorThree(%#v) returned %#v, expected %#v", second, res, expected)
	}
}

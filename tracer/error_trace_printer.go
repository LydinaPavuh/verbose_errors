package tracer

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var defaultOpts = FormatOpts{
	WithTrace:         true,
	WithUntracedWraps: true,
}

// FormatOpts format options for printing
type FormatOpts struct {
	WithTrace         bool // add trace on output
	WithCaller        bool // add caller (last trace line)
	WithUntracedWraps bool // add untraced errors between tracers
}

type frameList struct {
	Frame runtime.Frame
	Next  *frameList
}

func (fp *frameList) Format() string {
	res := fmt.Sprintln(fp.Frame.Function)
	res += fmt.Sprintf("	%s:%d\n", fp.Frame.File, fp.Frame.Line)
	return res
}

type StackPrinter struct {
	Head        *frameList
	Tail        *frameList
	NextPrinter *StackPrinter

	tailErr error
	wraps   []error
}

func PrintTrace(err error) string {
	return NewTracePrinter(err).Print()
}

func PrintTraceWithOpts(err error, opts *FormatOpts) string {
	return NewTracePrinter(err).PrintWithOpts(opts)
}

func NewTracePrinter(err error) *StackPrinter {

	root := &StackPrinter{tailErr: err}

	if tErr, ok := err.(ErrorTracer); ok {
		root.addTrace(tErr.Trace())
	}

	nexTracer, wrappedBy := unwrapTracer(err, []error{})
	root.wraps = wrappedBy
	if nexTracer != nil {
		root.NextPrinter = NewTracePrinter(nexTracer)
	}
	return root
}

func (p *StackPrinter) Print() string {
	return p.print(&defaultOpts)
}

func (p *StackPrinter) PrintWithOpts(opts *FormatOpts) string {
	return p.print(opts)
}

func (p *StackPrinter) print(opts *FormatOpts) string {
	res := ""
	if p.NextPrinter != nil {
		res += p.NextPrinter.print(opts)
	}

	if p.tailErr == nil {
		return res
	}
	return res + p.formatNode(opts)
}

func (p *StackPrinter) formatNode(opts *FormatOpts) string {
	res := strings.Repeat("====", 10) + "\n"
	res += p.formatError()

	if opts.WithUntracedWraps {
		res += p.formatWrappedErrors()
	}
	if opts.WithCaller {
		res += p.formatCaller()
	}
	if opts.WithTrace {
		res += p.formatTrace()
	}
	return res
}

func (p *StackPrinter) formatCaller() string {
	res := ""
	if p.Tail != nil {
		frame := p.Tail.Frame
		res += fmt.Sprintf("CALLER: [%s] %s:%d\n", frame.Function, frame.File, frame.Line)
	}
	return res
}
func (p *StackPrinter) formatTrace() string {
	res := ""
	if p.Tail != nil {
		res += fmt.Sprintln("TRACEBACK:")
		last := p.Head
		for last != nil {
			res += p.formatFrame(&last.Frame)
			last = last.Next
		}
	}
	return res
}

func (p *StackPrinter) formatFrame(frame *runtime.Frame) string {
	res := fmt.Sprintln("\t", frame.Function)
	res += fmt.Sprintf("\t\t %s:%d\n", frame.File, frame.Line)
	return res
}

func (p *StackPrinter) formatError() string {
	if p.tailErr == nil {
		return ""
	}
	return fmt.Sprintf("ERROR: %s\n", p.tailErr)
}

func (p *StackPrinter) formatWrappedErrors() string {
	if len(p.wraps) == 0 {
		return ""
	}
	res := fmt.Sprintln("UNTRACED WRAPS:")

	for i, err := range p.wraps {
		if err != nil {
			res += fmt.Sprintf("\t %d. %s\n", i, err.Error())
		}
	}
	return res
}

func (p *StackPrinter) addTrace(frames *runtime.Frames) {
	if frames == nil {
		return
	}
	frame, more := frames.Next()
	p.Head = &frameList{frame, p.Head}
	p.Tail = p.Head

	for more {
		p.Head = &frameList{Next: p.Head}
		p.Head.Frame, more = frames.Next()
	}
}

func unwrapTracer(e error, wrappedBy []error) (ErrorTracer, []error) {
	nextErr := errors.Unwrap(e)
	if nextErr == nil {
		return nil, wrappedBy
	}

	tErr, ok := nextErr.(ErrorTracer)
	if ok {
		return tErr, wrappedBy
	}

	return unwrapTracer(nextErr, append(wrappedBy, nextErr))
}

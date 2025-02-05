package errs

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type FormatStackFn func(stack Stack) string
type FormatFrameFn func(frame StackFrame) string
type FormatWrapFn func(frame StackFrame) error

var FormatError FormatStackFn = func(frame Stack) string {
	return fmt.Sprintf("%s ≡{%d}", frame[0].Error(), len(frame))
}

var FormatFrame FormatFrameFn = func(frame StackFrame) string {
	return fmt.Sprintf("%s:%d %s", filepath.Base(frame.File), frame.Line, frame.Err.Error())
}

var FormatStack FormatStackFn = func(stack Stack) string {
	errs := make([]string, len(stack)+1)
	errs[0] = stack.Error()

	for i := range stack {
		prefix := "│"
		for p := 1; p < i; p++ {
			prefix += " "
		}
		if i > 0 {
			prefix += "└"
		}
		if i == len(stack)-1 {
			prefix += "─"
		} else {
			prefix += "┬"
		}
		errs[i+1] = prefix + " " + stack[i].Dump()
	}
	return strings.Join(errs, "\n")
}

var FormatWrap FormatWrapFn = func(frame StackFrame) error {
	var stack Stack
	if !errors.As(frame.Err, &stack) {
		return fmt.Errorf("%s.%s » %w", filepath.Base(frame.Package), frame.Func, frame.Err)
	}
	return fmt.Errorf("%s.%s", filepath.Base(frame.Package), frame.Func)
}

type Stack []StackFrame

func (s Stack) Error() string {
	if len(s) == 0 {
		return "empty error stack"
	}
	if len(s) > 1 {
		return FormatError(s)
	}
	return s[0].Error()
}

func (s Stack) Unwrap() []error {
	errs := make([]error, len(s))
	for i := range s {
		errs[i] = s[i].Err
	}
	return errs
}

type StackFrame struct {
	Package string
	Func    string
	File    string
	Line    int
	Err     error
}

func (s StackFrame) Error() string {
	return s.Err.Error()
}

func (s StackFrame) Dump() string {
	return FormatFrame(s)
}

func newStackFrame(err error) StackFrame {
	pc, filename, line, _ := runtime.Caller(2)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]
	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return StackFrame{
		Package: packageName,
		Func:    funcName,
		File:    filename,
		Line:    line,
		Err:     err,
	}
}

func Wrap(err error) error {
	var eStack Stack
	frame := newStackFrame(err)
	frame.Err = FormatWrap(frame)
	errors.As(err, &eStack)
	return Stack(append([]StackFrame{frame}, eStack...))
}

func Wrapf(err error, format string, a ...any) error {
	text := fmt.Sprintf(format, a...)
	var eStack Stack
	if errors.As(err, &eStack) {
		return Stack(append([]StackFrame{newStackFrame(errors.New(text))}, eStack...))
	}
	return Stack([]StackFrame{newStackFrame(err)})
}

func Dump(err error) string {
	var eStack Stack
	if errors.As(err, &eStack) {
		return FormatStack(eStack)
	}
	return err.Error()
}

func Root(err error) string {
	var eStack Stack
	if errors.As(err, &eStack) && len(eStack) > 0 {
		return eStack[len(eStack)-1].Error()
	}
	return err.Error()
}

func Errorf(format string, a ...any) error {
	return Stack([]StackFrame{newStackFrame(fmt.Errorf(format, a...))})
}

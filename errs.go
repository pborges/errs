package errs

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type Stack []StackFrame

func (s Stack) Root() error {
	if len(s) == 0 {
		return nil
	}
	return s[len(s)-1].Err
}

func (s Stack) Error() string {
	if len(s) == 0 {
		return "empty error stack"
	}
	return s.Root().Error()
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
	msg := fmt.Sprintf("%s:%d (%s.%s)",
		filepath.Base(s.File),
		s.Line,
		filepath.Base(s.Package),
		s.Func,
	)
	if s.Err != nil {
		msg += " " + s.Err.Error()
	}
	return msg
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

// Wrap pushes the current error into the stack and places a new error on top with the default text provided by FormatWrap
func Wrap(err error) error {
	var eStack Stack
	if errors.As(err, &eStack) {
		return Stack(append([]StackFrame{newStackFrame(nil)}, eStack...))
	}
	return Stack([]StackFrame{newStackFrame(err)})
}

// Wrapf pushes the current error into the stack and places a new error on top with the formated text,
// it does not use the traditional error wrapping to do so
func Wrapf(err error, format string, a ...any) error {
	var eStack Stack
	if errors.As(err, &eStack) {
		return Stack(
			append(
				[]StackFrame{newStackFrame(fmt.Errorf(format, a...))},
				eStack...,
			),
		)
	}
	return Stack([]StackFrame{newStackFrame(fmt.Errorf(format+" » %s", append(a, err.Error())...))})
}

func Dump(err error) string {
	var stack Stack
	if errors.As(err, &stack) && len(stack) > 1 {
		errs := make([]string, len(stack)+1)
		errs[0] = stack.Error()
		for i := 0; i < len(stack); i++ {
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
			errs[i+1] = fmt.Sprintf("%s %s", prefix, stack[i].Error())
		}
		return strings.Join(errs, "\n")

	}
	return err.Error()
}

package errs

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type Stack []StackFrame

func (s Stack) Error() string {
	if len(s) == 0 {
		return "empty error stack"
	}
	return s[0].Error()
}

func (s Stack) Unwrap() []error {
	errs := make([]error, len(s))
	for i := range s {
		errs[i] = s[i].Cause
	}
	return errs
}

type StackFrame struct {
	Package string
	Func    string
	File    string
	Line    int
	Message string
	Cause   error
}

func (s StackFrame) Error() string {
	return s.Cause.Error()
}

func newStackFrame(err error, message string) StackFrame {
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
		Message: message,
		Cause:   err,
	}
}

// Wrap pushes the current error into the stack and places a new error on top with the default text provided by FormatWrap
func Wrap(err error) error {
	var eStack Stack
	errors.As(err, &eStack)
	return Stack(append([]StackFrame{newStackFrame(err, "")}, eStack...))
}

// Wrapf pushes the current error into the stack and places a new error on top with the formated text,
// it does not use the traditional error wrapping to do so
func Wrapf(err error, format string, a ...any) error {
	var eStack Stack
	if errors.As(err, &eStack) {
		return Stack(
			append(
				[]StackFrame{newStackFrame(err, fmt.Sprintf(format, a...))},
				eStack...,
			),
		)
	}
	return Stack([]StackFrame{newStackFrame(err, fmt.Sprintf(format+" » %s", append(a, err.Error())...))})
}

// Errorf wraps an error using the traditional fmt.Errorf method in an errs.Sack error
//func Errorf(format string, a ...any) error {
//	return Stack([]StackFrame{newStackFrame(fmt.Errorf(format, a...))})
//}

func Dump(err error) string {
	var stack Stack
	if errors.As(err, &stack) {
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
			errs[i+1] = fmt.Sprintf("%s %s:%d %s.%s %s",
				prefix,
				filepath.Base(stack[i].File),
				stack[i].Line,
				filepath.Base(stack[i].Package),
				stack[i].Func,
				stack[i].Message,
			)
		}
		return strings.Join(errs, "\n")
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

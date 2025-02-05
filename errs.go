package errs

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

var rootLast = true

type Stack []StackFrame

func (s Stack) Error() string {
	if len(s) == 0 {
		return "empty error stack"
	}
	if s[0].Err == nil {
		return ""
	}
	return s[0].Err.Error()
}

func (s Stack) Unwrap() []error {
	errs := make([]error, 0, len(s))
	for i := range s {
		if s[i].Err != nil {
			errs = append(errs, s[i].Err)
		}
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
		return append(eStack, newStackFrame(nil))
	}
	return Stack([]StackFrame{newStackFrame(err)})
}

// Wrapf pushes the current error into the stack and places a new error on top with the formated text,
// it does not use the traditional error wrapping to do so
func Wrapf(err error, format string, a ...any) error {
	var eStack Stack
	if errors.As(err, &eStack) {
		return append(eStack, newStackFrame(fmt.Errorf(format, a...)))
	}
	return Stack([]StackFrame{newStackFrame(fmt.Errorf(format+" » %w", append(a, err)...))})
}

func Detailed(err error) string {
	var stack Stack
	if errors.As(err, &stack) && len(stack) > 1 {
		lines := make([]string, len(stack)+1)
		lines[0] = stack.Error()
		for i := 0; i < len(stack); i++ {
			n := i
			if rootLast {
				n = len(stack) - 1 - i
			}
			prefix := "│"
			for p := 1; p < n; p++ {
				prefix += " "
			}
			if n > 0 {
				prefix += "└"
			}
			if n == len(stack)-1 {
				prefix += "─"
			} else {
				prefix += "┬"
			}
			lines[n+1] = fmt.Sprintf("%s %s", prefix, stack[i].Error())
		}
		return strings.Join(lines, "\n")

	}
	return err.Error()
}

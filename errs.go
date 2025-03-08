package errs

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type Stack struct {
	error
	Trace []StackFrame
}

func (s Stack) Unwrap() []error {
	errs := make([]error, len(s.Trace)+1)
	errs[0] = s.error
	for i := range s.Trace {
		errs[i+1] = s.Trace[i].error
	}
	return errs
}

type StackFrame struct {
	Package string
	Func    string
	File    string
	Line    int
	error
}

func (s StackFrame) String() string {
	msg := fmt.Sprintf("%s:%d (%s.%s)",
		filepath.Base(s.File),
		s.Line,
		filepath.Base(s.Package),
		s.Func,
	)
	if s.error != nil {
		msg += " " + transform(s.error)
	} else {
		msg += " ^"
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
		error:   err,
	}
}

// Push pushes err onto the stack if err is not nil
func Push(err error) error {
	if err == nil {
		return nil
	}
	var eStack Stack
	if !errors.As(err, &eStack) {
		eStack.error = err
		eStack.Trace = append(eStack.Trace, newStackFrame(err))
	} else {
		eStack.Trace = append(eStack.Trace, newStackFrame(nil))
	}
	return eStack
}

// Wrap pushes err and wrapped onto the stack if err is not nil
func Wrap(err error, wrapped error) error {
	if err == nil {
		return nil
	}
	var eStack Stack
	if errors.As(err, &eStack) {
		eStack.Trace = append(eStack.Trace, newStackFrame(wrapped))
	} else {
		eStack.Trace = append(eStack.Trace, newStackFrame(err), newStackFrame(wrapped))
	}
	eStack.error = wrapped

	return eStack
}

func Detailed(err error) string {
	if err == nil {
		return ""
	}
	var stack Stack
	if errors.As(err, &stack) && len(stack.Trace) > 1 {
		trace := stack.Trace
		var message []string
		for i, line := range strings.Split(transform(stack), "\n") {
			line = strings.TrimSpace(line)
			if i == 0 {
				message = append(message, line)
			} else {
				message = append(message, "│ "+line)
			}
		}
		return strings.Join(message, "\n") + "\n" + stepJoin(trace, StackFrame.String)
	}
	return err.Error()
}

func stepJoin[T any](input []T, stringer func(T) string) string {
	var results []string
	for i := 0; i < len(input); i++ {
		prefix := "│"
		for p := 1; p < i; p++ {
			prefix += " "
		}
		if i > 0 {
			prefix += "└"
		}
		if i == len(input)-1 {
			prefix += "─"
		} else {
			prefix += "┬"
		}
		output := stringer(input[len(input)-1-i])
		first := true
		for line := range strings.Lines(output) {
			line = strings.TrimSpace(line)
			if first {
				results = append(results, prefix+" "+line)
				first = false
			} else {
				multilinePrefix := "│"
				for p := 1; p < i; p++ {
					multilinePrefix += " "
				}
				results = append(results, multilinePrefix+"   "+line)
			}
		}
	}
	return strings.Join(results, "\n")
}

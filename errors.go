package goexer

import (
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	Message   string // Error message.
	Code      int    // Error code (usage is optional).
	Function  string // Function where error was created.
	File      string // File with code definition.
	Line      uint   // Line where error was created.
	Previous  *Error // Previous error.
	Original  error  // Original error if wrap was used for non Error objects.
	Container Container
}

func (e *Error) Error() string {
	return e.Message
}

// setLocation - record where is error was created.
func (e *Error) setLocation(depth int) {
	pc := make([]uintptr, 1)
	runtime.Callers(depth+2, pc) // Skip setLocation and Callers itself.
	frame, _ := runtime.CallersFrames(pc).Next()
	e.Function = frame.Function
	e.Line = uint(frame.Line)
	file := strings.Split(frame.File, "/")
	e.File = file[len(file)-1]
}

func (e *Error) Set(key string, value interface{}) {
	e.Container.Set(key, value)
}

func (e *Error) Get(key string) any {
	return e.Container.Get(key)
}

func (e *Error) GetRaw(key string) any {
	return e.Container.GetRaw(key)
}

func (e *Error) GetE(key string) (any, bool) {
	return e.Container.GetE(key)
}

func (e *Error) GetRawE(key string) (any, bool) {
	return e.Container.GetRawE(key)
}

func (e *Error) Cause() error {
	return e.Original
}

// Pretty string for error.
func (e *Error) PrettyError() string {
	return fmt.Sprintf("%s:%d %s %s", e.File, e.Line, e.Function, e.Message)
}

// Return slice of errors in correct sequence.
func (e *Error) Stack() []*Error {
	stack := []*Error{}

	err := e
	for err != nil {
		stack = append([]*Error{err}, stack...)
		err = err.Previous
	}

	return stack
}

// Return stack as pretty string.
func (e *Error) StackString() string {
	s := ""

	for _, err := range e.Stack() {
		s += err.PrettyError() + "\n"
	}

	return s
}

// Allow to print %#v (I need better method, because this brakes type name).
type ErrorWrap Error

func (e *ErrorWrap) Format() {}

// Format - implements extended fmt.Formatter.
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			fmt.Fprintf(s, "%s", e.StackString())

			return
		case s.Flag('#'):
			// Allow to use %#v for Error.
			fmt.Fprintf(s, "%#v", (*ErrorWrap)(e))

			return
		}

		fallthrough

	case 's':
		fmt.Fprintf(s, "%s", e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	default:
		fmt.Fprintf(s, "%%!%c(%T=%s)", verb, e, e.Error())
	}
}

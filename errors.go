package goexer

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

type Error struct {
	Message              string // Error message.
	Name                 string // Error name (kind). E.g. NotFound, ... Do not use long strings for better formatting.
	Function             string // Function where error was created.
	File                 string // File with code definition.
	Line                 uint   // Line where error was created.
	Previous             *Error // Previous error.
	Original             error  // Original error if wrap was used for non Error objects.
	Container            *Container
	ShowContainerItems   []string
	ShowContainerSize    bool
	ShowContainerAsZKeys bool // Show container items as keys in zerolog, instead of as message.
	AddTraceToError      bool // Add trace messages to error and fatal messages with error level ERROR.
}

// Additional options for New(), Wrap(), ...
type ErrorOpts struct {
	Name                 string     // Name (kind) of error. Do not use long strings for better formatting.
	Depth                int        // Depth of stack trace. Increase if need fetch data from previous frame.
	Container            *Container // Use existing container with prefilled data. Use nil for personal container for each error.
	ShowContainerSize    *bool      // Show container size in error message.
	ShowContainerItems   []string   // Show container items in error message (key and value).
	ShowContainerAsZKeys *bool      // Show container items as keys in zerolog, instead of as message.
	AddTraceToError      *bool      // Add trace messages to error and fatal messages with error level ERROR.
}

func (e *Error) Error() string {
	return e.OneLinePrettyError()
	// return e.Message
}

// setLocation - record where is error was created.
func (e *Error) setLocation(depth int) {
	pc := make([]uintptr, 1)
	runtime.Callers(depth+2, pc) // Skip setLocation and Callers itself.
	frame, _ := runtime.CallersFrames(pc).Next()
	e.Function = frame.Function
	e.Line = uint(frame.Line)
	// file := strings.Split(frame.File, "/")
	// e.File = file[len(file)-1]
	e.File = frame.File
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

// Pretty string for error in one line style.
func (e *Error) OneLinePrettyError() string {
	file := strings.Split(e.File, "/")

	s := fmt.Sprintf("%s: %s:%d %s(): '%s'", e.Name, file[len(file)-1], e.Line, e.Function, e.Message)

	if e.ShowContainerSize && !e.ShowContainerAsZKeys {
		s += fmt.Sprintf(" (cSize: %d)", e.Container.Size())
	}
	if e.ShowContainerItems != nil && !e.ShowContainerAsZKeys {
		s += " (Fields:"
		for _, i := range e.ShowContainerItems {
			s += fmt.Sprintf(" %s: %v;", i, e.Get(i))
		}
		s += ")"
	}

	return s
}

// Pretty string for error in one line style.
func (e *Error) MultiLinePrettyError() string {
	s := fmt.Sprintf("%s(): %s\n\t%s:%d\n", e.Function, e.Message, e.File, e.Line)
	if e.ShowContainerSize {
		s += fmt.Sprintf("\tContainer size: %d\n", e.Container.Size())
	}
	if e.ShowContainerItems != nil {
		s += "\tContainer fields:"
		for _, i := range e.ShowContainerItems {
			s += fmt.Sprintf(" %s: %v;", i, e.Get(i))
		}
		s += "\n"
	}

	return s
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
		s += err.MultiLinePrettyError()
	}

	return s
}

// Support for errors.Is().
// Return true if err.Name == e.Name .
func (e *Error) Is(err error) bool {
	if target, ok := err.(*Error); ok {
		return e.Name == target.Name
	}

	return false
}

// Support errors.Unwrap().
func (e *Error) Unwrap() error {
	return e.Previous
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

// Log - log error. Currently supports only zerolog. msg - Additional message. Could be empty.
func (e *Error) log(event *zerolog.Event, msg ...string) {
	// Use this only for trace and/or debug level.
	// newMsg := e.StackString()
	newMsg := ""
	for _, m := range msg {
		newMsg += " " + m
	}

	if e.ShowContainerSize && e.ShowContainerAsZKeys {
		event.Int("cSize", e.Container.Size())
	}
	if e.ShowContainerItems != nil && e.ShowContainerAsZKeys {
		for _, i := range e.ShowContainerItems {
			event.Str("field_"+i, fmt.Sprintf("%v", e.Get(i)))
		}
	}

	event.Str("error", e.OneLinePrettyError()).Msg(newMsg)
}

// Log with Error log level.
func (e *Error) LogError(msg ...string) {
	if zLog != nil {
		if e.AddTraceToError {
			e.LogTraceToEvent(zLog.Error())
		}
		e.log(zLog.Error(), msg...)
	}
}

// Log with fatal log level (with program exit!).
func (e *Error) LogFatal(msg ...string) {
	if zLog != nil {
		if e.AddTraceToError {
			e.LogTraceToEvent(zLog.Error())
		}
		e.log(zLog.Fatal(), msg...)
	}
}

// Log trace to event.
func (e *Error) LogTraceToEvent(event *zerolog.Event, msg ...string) {
	if zLog != nil {
		newMsg := ""
		for _, m := range msg {
			newMsg += " " + m
		}
		event.Msg(newMsg + "\n" + e.StackString())
	}
}

// Log with trace log level.
func (e *Error) LogTrace(msg ...string) {
	e.LogTraceToEvent(zLog.Trace())
}

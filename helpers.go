package goexer

import (
	"fmt"
	"log"
	"reflect"

	"github.com/rs/zerolog"
)

// Check if error is *goexer.Error.
func IsGoexerError(err error) bool {
	return reflect.TypeOf(err).String() == ErrorTypeString
}

// New - create new Error.
func newError(depth int, msg string, op ErrorOpts) *Error {
	err := Error{}
	opts := DefaultErrorOpts

	switch {
	case op.Name != "":
		opts.Name = op.Name
	case op.Container != nil:
		opts.Container = op.Container
	case op.ShowContainerItems != nil:
		opts.ShowContainerItems = op.ShowContainerItems
	case op.ShowContainerSize != nil:
		opts.ShowContainerSize = op.ShowContainerSize
	case op.ShowContainerAsZKeys != nil:
		opts.ShowContainerAsZKeys = op.ShowContainerAsZKeys
	}

	err.setLocation(depth + opts.Depth) // This error.
	if opts.Container == nil {
		err.Container = NewContainer()
	} else {
		err.Container = opts.Container
	}
	err.Message = msg
	err.Name = opts.Name

	err.ShowContainerItems = opts.ShowContainerItems
	if opts.ShowContainerAsZKeys != nil {
		err.ShowContainerAsZKeys = *opts.ShowContainerAsZKeys
	}
	if opts.ShowContainerSize != nil {
		err.ShowContainerSize = *opts.ShowContainerSize
	}

	return &err
}

// New - create new Error.
func New(msg string, args ...ErrorOpts) *Error {
	if len(args) > 1 {
		fatal(New("Only one or zero ErrorOpts could be passed to New()"), "Only one or zero ErrorOpts could be passed to New()")
	}

	opts := DefaultErrorOpts

	if len(args) == 1 {
		opts = args[0]
	}

	ee := newError(2+opts.Depth, msg, opts)

	return ee
}

func fatal(err *Error, msg string) {
	switch {
	case zLog != nil:
		zLog.Fatal().
			Uint("line", err.Line).
			Str("function", err.Function).
			Str("file", err.File).
			Msg(msg)
	case bLog != nil:
		bLog.Fatalf(
			msg+" line: %d func: %s file: %s\n",
			err.Line,
			err.Function,
			err.File,
		)
	default:
		// Last chance.
		panic("Incorrect wrap usage. Previous error should not be nil.")
	}
}

// Wrapp old error to the new one.
func Wrap(prev error, msg string, args ...ErrorOpts) *Error {
	if len(args) > 1 {
		fatal(New("Only one or zero ErrorOpts could be passed to Wrap()"), "Only one or zero ErrorOpts could be passed to Wrap()")
	}

	opts := DefaultErrorOpts

	if len(args) == 1 {
		opts = args[0]
	}

	err := newError(2, msg, opts) // Current error stack.

	if prev == nil {
		fatal(err, "goexer.Error.Wrap: Incorrect wrap usage. Previous error should not be nil.")
	}

	if IsGoexerError(prev) {
		err.Name = ToError(prev).Name
	}

	if !IsGoexerError(prev) {
		errPrev := newError(3, prev.Error(), opts) // Previous error stack.
		errPrev.Original = prev
		err.Previous = errPrev
		err.Original = prev
	} else {
		err.Previous = ToError(prev)
		if err.Previous.Original == nil {
			err.Original = err.Previous
		} else {
			err.Original = err.Previous.Original
		}
	}

	return err
}

// Convert any error to Error.
func ToError(err error) *Error {
	ee, ok := err.(*Error)
	if !ok {
		return newError(3, err.Error(), DefaultErrorOpts) // Previous error stack.
	}

	return ee
}

// Formatted wrap.
func Wrapf(prev error, format string, args ...interface{}) *Error {
	return Wrap(prev, fmt.Sprintf(format, args...))
}

func Cause(err error) error {
	if IsGoexerError(err) {
		return ToError(err).Original
	}

	return err
}

// SetOpts - set DefaultOptions.
func SetDefaultOpts(opts ErrorOpts) {
	DefaultErrorOpts = opts
}

// Set global zerolog logger.
func SetZLog(log *zerolog.Logger) {
	zLog = log
}

// Set global basic logger.
func SetBLog(log *log.Logger) {
	bLog = log
}

// Check error and log fatal message if not nil.
func CheckErr(err error) {
	if err != nil {
		ToError(err).LogFatal()
	}
}

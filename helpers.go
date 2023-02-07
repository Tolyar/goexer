package goexer

import (
	"fmt"
	"reflect"
)

// New - create new Error.
func newError(depth int, msg string) *Error {
	err := Error{}
	err.setLocation(depth) // This error.
	err.Container = NewContainer()
	err.Message = msg

	return &err
}

// New - create new Error.
func New(msg string) *Error {
	ee := newError(2, msg)

	return ee
}

// Formatted new.
func Newf(format string, args ...interface{}) *Error {
	return New(fmt.Sprintf(format, args...))
}

func fatal(err *Error, msg string) {
	//nolint:gocritic
	if zLog != nil {
		zLog.Fatal().
			Uint("line", err.Line).
			Str("function", err.Function).
			Str("file", err.File).
			Msg(msg)
	} else if bLog != nil {
		bLog.Fatalf(
			msg+" line: %d func: %s file: %s\n",
			err.Line,
			err.Function,
			err.File,
		)
	} else {
		// Last chance.
		panic("Incorrect wrap usage. Previous error should not be nil.")
	}
}

// Wrapp old error to the new one.
func Wrap(prev error, msg string) *Error {
	err := newError(2, msg) // Current error stack.

	if prev == nil {
		fatal(err, "goexer.Error.Wrap: Incorrect wrap usage. Previous error should not be nil.")
	}

	if reflect.TypeOf(prev).String() != "*goexer.Error" {
		errPrev := newError(3, prev.Error()) // Previous error stack.
		errPrev.Original = prev
		err.Previous = errPrev
		err.Original = prev
	} else if p, ok := prev.(*Error); ok {
		err.Previous = p
		if p.Original == nil {
			err.Original = p
		} else {
			err.Original = p.Original
		}
	} else {
		fatal(err, "Can't convert prev to *Error")
	}

	return err
}

// Formatted wrap.
func Warpf(prev error, format string, args ...interface{}) *Error {
	return Wrap(prev, fmt.Sprintf(format, args...))
}

func Cause(err error) error {
	if reflect.TypeOf(err).String() == "*goexer.Error" {
		ee, ok := err.(*Error)
		if !ok {
			fatal(newError(2, "goexer.Cause"), "goexer.Cause can't convert interface{}")
		}

		return ee.Original
	}

	return err
}

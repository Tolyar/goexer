package goexer_test

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/Tolyar/goexer"
)

func TestNew(t *testing.T) {
	t.Parallel()

	err := goexer.New("Test error")
	//nolint:dogsled
	_, currentFile, _, _ := runtime.Caller(0)

	if reflect.TypeOf(err.Container).String() != "*goexer.Container" {
		t.Errorf("Want container type 'goexec.Container' but got '%T'", err.Container)
	}
	if err.Line == 0 {
		t.Errorf("Want line > 0, got '%d'", err.Line)
	}

	if err.Function == "TestNew" {
		t.Errorf("Function name should be TestNew, got %s", err.Function)
	}

	if err.File != currentFile {
		t.Errorf("File name should be %s, got %s", currentFile, err.File)
	}
}

func TestWrap(t *testing.T) {
	t.Parallel()

	//nolint
	eo := errors.New("original")
	ep := goexer.Wrap(eo, "previous")
	err := goexer.Wrap(ep, "current")

	// Check current error.
	if err.Line == 0 {
		t.Errorf("Want line > 0, got '%d'", err.Line)
	}

	if err.Function != "github.com/Tolyar/goexer_test.TestWrap" {
		t.Errorf("Function name should be github.com/Tolyar/goexer_test.TestWrap, got %s", err.Function)
	}

	//nolint:dogsled
	_, currentFile, _, _ := runtime.Caller(0)
	//nolint:dogsled
	_, prevFile, _, _ := runtime.Caller(1)

	if err.File != currentFile {
		t.Errorf("File name should be %s, got %s", currentFile, err.File)
	}
	if err.Previous == nil {
		t.Error("err.Previous should not be nil")
	} else {
		err = err.Previous
	}

	// Check previous error.
	if reflect.TypeOf(err).String() != "*goexer.Error" {
		t.Errorf("Want Error type, got '%T'", err)
	}
	if err.Line == 0 {
		t.Errorf("Want line > 0, got '%d'", err.Line)
	}

	if err.Function != "github.com/Tolyar/goexer_test.TestWrap" {
		t.Errorf("Function name should be github.com/Tolyar/goexer_test.TestWrap, got %s", err.Function)
	}

	if err.File != currentFile {
		t.Errorf("File name should be %s, got %s", currentFile, err.File)
	}
	if err.Previous == nil {
		t.Error("err.Previous should not be nil")
	} else {
		err = err.Previous
	}

	// Check original error.
	if err.Line == 0 {
		t.Errorf("Want line > 0, got '%d'", err.Line)
	}

	if err.Function != "testing.tRunner" {
		t.Errorf("Function name should be testing.tRunner, got %s", err.Function)
	}

	// //nolint:dogsled
	// _, file, _, _ := runtime.Caller(0)
	// fs := strings.Split(file, "/")
	// currentFile = fs[len(fs)-1]

	if err.File != prevFile {
		t.Errorf("File name should be %s, got %s", prevFile, err.File)
	}

	//nolint:goconst
	if err.Original == nil || err.Original.Error() != "original" {
		t.Errorf("Incorrect original error. Wants not nil with message 'original', got '%v'", err.Original)
	}
}

func TestTypeCause(t *testing.T) {
	t.Parallel()

	//nolint:goerr113
	err := goexer.Wrap(goexer.Wrap(errors.New("original"), "previous"), "current")

	if err.Cause().Error() != "original" {
		t.Errorf("Want %s, got %s", "original", err.Cause().Error())
	}
}

func TestFuncCause(t *testing.T) {
	t.Parallel()

	//nolint:goerr113
	err := goexer.Wrap(goexer.Wrap(errors.New("original"), "previous"), "current")
	if goexer.Cause(err).Error() != "original" {
		t.Errorf("Want %s, got %s", "original", err.Cause().Error())
	}

	//nolint:goerr113
	er := errors.New("test")
	if goexer.Cause(er).Error() != "test" {
		t.Errorf("Want %s, got %s", "test", goexer.Cause(er).Error())
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	pc, file, line, _ := runtime.Caller(0)
	//nolint:goerr113
	err := goexer.Wrap(goexer.Wrap(errors.New("orig"), "previous"), "current")

	str := fmt.Sprintf("%s", err)

	fs := strings.Split(file, "/")
	currentFile := fs[len(fs)-1]
	// BaseError: errors_test.go:139 github.com/Tolyar/goexer_test.TestFormat(): 'current'
	msg := fmt.Sprintf("BaseError: %s:%d %s(): 'current'", currentFile, line+2, runtime.FuncForPC(pc).Name())
	if str != msg {
		t.Errorf("%%s format: want '%s', got '%s'", msg, str)
	}

	str = fmt.Sprintf("%v", err)
	if str != msg {
		t.Errorf("%%v format: want '%s', got '%s'", msg, str)
	}

	str = fmt.Sprintf("%+v", err)
	want := err.StackString()
	if str != want {
		t.Errorf("%%+v format: want '%s', got '%s'", want, str)
	}

	want = `&goexer.ErrorWrap{Message:"current", Name:"BaseError", Function:"github.com/Tolyar/goexer_test.TestFormat"`
	str = fmt.Sprintf("%#v", err)[:len(want)]
	if strings.Compare(want, str) != 0 {
		t.Errorf("%%#v format:\n want '%s'\n  got '%s'", want, str)
	}
}

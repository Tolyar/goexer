package main

import (
	"errors"
	"os"

	"github.com/Tolyar/goexer"
	"github.com/rs/zerolog"
)

func ErrNotFound(msg string) *goexer.Error {
	return goexer.New(msg, goexer.ErrorOpts{Depth: 1, Name: "ErrNotFound"})
}

func main() {
	zlog := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	goexer.SetZLog(&zlog)

	t := true
	c := goexer.NewContainer().Set("testInt", 10).Set("testString", "zzzz")
	goexer.SetDefaultOpts(
		goexer.ErrorOpts{
			Container:            c,
			ShowContainerSize:    &t,
			ShowContainerItems:   []string{"testInt", "testString", "testBool"},
			ShowContainerAsZKeys: &t,
		},
	)

	//nolint:goerr113
	e := errors.New("stderr")
	err := goexer.Wrap(e, "New() w/o opts")

	err.ShowContainerAsZKeys = false
	err.LogError("message1")
	err.ShowContainerAsZKeys = true
	err.LogError("message2")
	err.LogTrace("message")

	err.Set("testBool", true)
	ErrNotFound("NotFound").LogTrace()

	// Check Is().
	zlog.Info().Msgf("errors.Is() - %v (should be true)", errors.Is(ErrNotFound("test"), goexer.New("", goexer.ErrorOpts{Name: "ErrNotFound"})))

	err.LogFatal()
	err.LogError("This message should not be printed after log.Fatal()")
}

# goexer - Golang extended errors package

**WARNING: API of this package is still unstable. Please wait v1.0.0 before usage.**

I like Golang errors extenders but couldn't find a library that is ideal for me.

E.g.

* https://github.com/pkg/errors - my favorite package, but doesn't allow to easily extend errors with custom fields.
* https://upspin.io/errors - supports custom fields, but they are fixed inside library.
* https://github.com/juju/errors/ - Very useful package, but doesn't support custom fields.

In this project I will try to implement a flexible error type. I still don't know how it should work, but I have several ideas. So, I do not recommend using this library before v1.0.0 will be done.

# How it works

Each error has an embedded storage container and stack of errors. At any time you can get any error from stack and check its type, message and where it was raised (function name, file path, line number). You can Wrap() any errors interface compatible errors. Goexer does not full compatible with errors package and can't replace it. Also some functions may works in different style than functions from errors.
Each error can be created with ErrorOpts. By default will be used DefaultErrorOpts. Goexer has some helpers for usage with zerolog. You need set zerolog's instance for usage it.

# Types

## ErrorOpts - options for errors creations.
```go
type ErrorOpts struct {
	Name                 string     // Name (kind) of error. Do not use long strings for better formatting.
	Depth                int        // Depth of stack trace. Increase if need fetch data from previous frame.
	Container            *Container // Use existing container with prefilled data. Use nil for personal container for each error.
	ShowContainerSize    *bool      // Show container size in error message.
	ShowContainerItems   []string   // Show container items in error message (key and value).
	ShowContainerAsZKeys *bool      // Show container items as keys in zerolog, instead of as message.
}
```

## Container - storage for additional fields
```go
type Container struct {
	items map[string]Item
}
```



## Item - type for additional fields. Using inside Container.
```go
type Item struct {
	Value interface{}
	Type  string // Just for better string representation.
	Name  string // Just for better string representation.
}
```

## Error - main type for this package.

In most cases you will works only with this type. This type is compatible with error interface.

```go
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
}
```

# Examples

Set zerolog's instance.

```go
zlog := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
goexer.SetZLog(&zlog)
```

Set default for all errors in your project.
```go
t := true
c := goexer.NewContainer().Set("testInt", 10).Set("testString", "zzzz")
goexer.SetDefaultOpts(
    goexer.ErrorOpts{
        Container:            c, // Set one container for all errors.
        ShowContainerSize:    &t, // Show size of container in errors.
        ShowContainerItems:   []string{"testInt", "testString", "testBool"}, // Set fields than sholud be printed in error from container.
        ShowContainerAsZKeys: &t, // Print field as separate keys in zerolog instead of set they as message.
    },
)
```

Wrap standard error
```go
e := errors.New("stderr")
err := goexer.Wrap(e, "New() w/o opts")
```

Print error via zerolog Error() event with info about container in message.
```go
err.ShowContainerAsZKeys = false
err.LogError("message1")
```

Print error as error and as trace with fields as keys in zerolog events.
```go
err.ShowContainerAsZKeys = true
err.LogError("message2")
err.LogTrace("message")
```

The correct way to create a new error. Use Depth: 1 for getting original .
```go
func ErrNotFound(msg string) *goexer.Error {
	return goexer.New(msg, goexer.ErrorOpts{Depth: 1})
}

err.Set("testBool", true)
ErrNotFound("NotFound").LogTrace()
err.LogFatal()
err.LogError("This message should not be printed after log.Fatal()")
```

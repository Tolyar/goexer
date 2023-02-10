package goexer

import (
	"log"

	"github.com/rs/zerolog"
)

var (
	zLog *zerolog.Logger // Zerolog logger.
	bLog *log.Logger     // Basic logger ("log" package).
)

const (
	BaseErrorName   = "BaseError"
	ErrorTypeString = "*goexer.Error"
)

var DefaultErrorOpts ErrorOpts = ErrorOpts{
	Name:                 BaseErrorName,
	Depth:                0,
	Container:            nil,
	ShowContainerAsZKeys: nil,
	ShowContainerSize:    nil,
	ShowContainerItems:   nil,
}

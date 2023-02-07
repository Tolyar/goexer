package goexer

import (
	"log"

	"github.com/rs/zerolog"
)

var (
	zLog *zerolog.Logger // Zerolog logger.
	bLog *log.Logger     // Basic logger ("log" package).
)

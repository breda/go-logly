package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(service string) *zerolog.Logger {
	logger := zerolog.
		New(os.Stdout).
		With().
		Str("service", service).
		Timestamp().
		Logger()

	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	return &logger
}

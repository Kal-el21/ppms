package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init(env string) {
	if env == "production" {
		Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		Log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	}
}

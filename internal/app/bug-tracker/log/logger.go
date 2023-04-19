package log

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

//go:generate mockgen -source=logger.go -destination=mocks/logger.go

type Log interface {
	Info(msg string)
	Infof(msg string, v ...interface{})
	Fatal(err error)
	Error(err error)
	Internal() *zerolog.Logger
}

type Logger struct {
	log *zerolog.Logger
}

func New(out io.Writer) *Logger {
	logger := zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()

	return &Logger{
		log: &logger,
	}
}

func (l *Logger) Info(msg string) {
	l.log.Info().Msg(msg)
}

func (l *Logger) Infof(msg string, v ...interface{}) {
	l.log.Info().Msgf(msg, v...)
}

func (l *Logger) Error(err error) {
	l.log.Error().Err(err).Msg("")
}

func (l *Logger) Fatal(err error) {
	l.log.Fatal().Err(err).Msg("")
}

func (l *Logger) Internal() *zerolog.Logger {
	return l.log
}

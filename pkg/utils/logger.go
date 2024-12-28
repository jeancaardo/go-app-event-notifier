package utils

import (
	logS "log"
	"os"
	"reflect"
	"time"

	"github.com/getsentry/sentry-go"
)

const (
	sentryOption = "SentryOption"
	logOption    = "LogOption"
)

type SentryOption struct {
	Dsn         string
	Environment string
	Release     string
	FlushTime   time.Duration
	Debug       bool
}

type LogOption struct {
	Debug  bool
	Logger *logS.Logger
}

type Logger interface {
	CatchMessage(msg string) string
	CatchError(err error) error
	DebugMessage(msg string) string
	DebugError(err error) error
}
type log struct {
	logs []Logger
}

func New(loggers ...interface{}) Logger {
	l := log{}
	for _, ls := range loggers {
		l.add(ls)
	}
	return &l
}

func (l *log) CatchMessage(msg string) string {
	for _, log := range l.logs {
		log.CatchMessage(msg)
	}
	return msg
}

func (l *log) CatchError(msg error) error {
	for _, log := range l.logs {
		_ = log.CatchError(msg)
	}
	return msg
}

func (l *log) DebugMessage(msg string) string {
	for _, log := range l.logs {
		log.DebugMessage(msg)
	}

	return msg
}

func (l *log) DebugError(msg error) error {
	for _, log := range l.logs {
		_ = log.DebugError(msg)
	}
	return msg
}

func (l *log) add(log interface{}) {
	switch reflect.TypeOf(log).Name() {
	case sentryOption:
		so := log.(SentryOption)

		err := sentry.Init(sentry.ClientOptions{
			Dsn:         so.Dsn,
			Environment: so.Environment,
			Release:     so.Release,
		})

		l.logs = append(l.logs, &sentryLogger{err, so.Debug})

		defer sentry.Flush(so.FlushTime)
	case logOption:
		lo := log.(LogOption)

		if lo.Logger == nil {
			lo.Logger = logS.New(os.Stdout, "", logS.LstdFlags|logS.Lshortfile)
		}
		l.logs = append(l.logs, &logLogger{nil, lo.Debug, lo.Logger})

	}
}

type sentryLogger struct {
	err   error
	debug bool
}

func (sl *sentryLogger) CatchMessage(msg string) string {
	if msg != "" {
		sentry.CaptureMessage(msg)
	}
	return msg
}

func (sl *sentryLogger) CatchError(err error) error {
	if err != nil {
		sentry.CaptureException(err)
	}
	return err
}

func (sl *sentryLogger) DebugMessage(msg string) string {
	if msg != "" && sl.debug {
		sentry.CaptureMessage(msg)
	}
	return msg
}

func (sl *sentryLogger) DebugError(err error) error {
	if err != nil && sl.debug {
		sentry.CaptureException(err)
	}
	return err
}

type logLogger struct {
	err    error
	debug  bool
	logger *logS.Logger
}

func (sl *logLogger) CatchMessage(msg string) string {
	if msg != "" {
		sl.logger.Output(3, msg)
	}
	return msg
}

func (sl *logLogger) CatchError(err error) error {
	if err != nil {
		sl.logger.Output(3, err.Error())
	}
	return err
}

func (sl *logLogger) DebugMessage(msg string) string {
	if msg != "" && sl.debug {
		sl.logger.Output(3, msg)
	}
	return msg
}

func (sl *logLogger) DebugError(err error) error {
	if err != nil && sl.debug {
		sl.logger.Output(3, err.Error())
	}
	return err
}

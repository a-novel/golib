package anovelmonitor

import "io"

// partialLogger implements non-agnostic methods, so agnostic logs can be automatically added to form a fully
// functional Logger.
type partialLogger interface {
	Fatal(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Info(msg string, args ...interface{})

	io.Writer
}

// When fed a partialLogger (with any implementation), loggerImpl completes it to form the final Logger.
type loggerImpl struct {
	// Logger non-agnostic methods.
	partialLogger
}

func (l *loggerImpl) FatalE(err error, msg string) {
	// If no message is present, just print the error.
	if msg == "" {
		l.Fatal("%v", err)
	}

	// Concatenate the message and the error.
	l.Fatal("%s: %v", msg, err)
}

func (l *loggerImpl) ErrorE(err error, msg string) {
	// If no message is present, just print the error.
	if msg == "" {
		l.Error("%v", err)
	}

	// Concatenate the message and the error.
	l.Error("%s: %v", msg, err)
}

// NewLogger creates a new Logger from a partialLogger.
func newLoggerFromImpl(p partialLogger) Logger {
	return &loggerImpl{p}
}

package loggers

import (
	"io"
	"os"
)

// Implements Console using the io.writer interface.
type stdLogger struct {
	stdout  io.Writer
	stdwarn io.Writer
	stderr  io.Writer
}

// Log implements the Console.Log interface.
func (logger *stdLogger) Log(level LogLevel, msg string) {
	// Don't log empty messages, do a no-op instead.
	if msg == "" {
		return
	}

	// Switch channel depending on the log level.
	switch level {
	case LogLevelInfo:
		_, _ = io.WriteString(logger.stdout, msg)
	case LogLevelWarning:
		_, _ = io.WriteString(logger.stdwarn, msg)
	case LogLevelError:
		_, _ = io.WriteString(logger.stderr, msg)
	case LogLevelFatal:
		// Fatal error use the same error channel as regular errors. The program is killed after logging.
		_, _ = io.WriteString(logger.stderr, msg)
		os.Exit(1)
	default:
		_, _ = io.WriteString(logger.stdout, msg)
	}
}

// NewSTDOut creates a new Console that writes to the standard outputs defined by the OS.
//
// LogLevelInfo and LogLevelWarning are written to os.Stdout, LogLevelError and LogLevelFatal are written to os.Stderr.
func NewSTDOut() Console {
	return &stdLogger{
		stdout:  os.Stdout,
		stdwarn: os.Stdout,
		stderr:  os.Stderr,
	}
}

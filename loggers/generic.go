package loggers

// LogLevel specify the importance of the log message. Some implementations may use different channels depending
// on the log level.
type LogLevel string

const (
	// LogLevelInfo is the lowest log level. It is used for general information messages.
	LogLevelInfo LogLevel = "INFO"
	// LogLevelWarning is used for messages that are not errors but may require attention.
	LogLevelWarning LogLevel = "WARNING"
	// LogLevelError is used for messages that indicate an error occurred.
	LogLevelError LogLevel = "ERROR"
	// LogLevelFatal is used for messages that indicate a fatal error occurred. A logger implementation should
	// automatically exit the program, or trigger a crash, after logging a message with this level.
	LogLevelFatal LogLevel = "FATAL"
)

// Console prints an output to a live terminal.
type Console interface {
	// Log prints the message to the console.
	Log(level LogLevel, msg string)
}

// JSON sends a serialized JSON representation of the log message.
type JSON interface {
	// Log sends a JSON representation of the log message.
	Log(level LogLevel, data interface{})
}

package formatters

import "github.com/a-novel/golib/loggers"

// LogContent represents a special value, that can automatically be rendered in both logger.Console and
// logger.JSON.
type LogContent interface {
	// RenderConsole returns a string representation of the log content, that can be printed in a logger.Console.
	RenderConsole() string
	// RenderJSON returns a serializable representation of the log content, that can be printed in a logger.JSON.
	RenderJSON() interface{}
}

// LogDynamicContent allows to render real-time output when the logger supports it. It is currently only available
// for logger.Console.
type LogDynamicContent interface {
	LogContent
	// RenderConsoleDynamic is a special method, that will be called periodically by the logger to update the content.
	// It is up to the dynamic implementation to keep an internal running state. This method will be called everytime
	// the message should be updated.
	//
	// This method must be ignored once StopRunning has been called.
	RenderConsoleDynamic(renderer func(msg string)) LogDynamicContent
	// StopRunning stops the dynamic render execution. After triggering it, every call to RenderConsoleDynamic must be
	// ignored.
	StopRunning() LogDynamicContent
}

// Formatter will handle a LogContent message and dispatch it to the appropriate logger.
type Formatter interface {
	// Log sends a content to the logger, with the specified log level.
	Log(content LogContent, level loggers.LogLevel)
}

// Default implementation of the Formatter interface, for logger.Console instances.
type consoleFormatterImpl struct {
	// The logger that will be used to print the message.
	logger loggers.Console
	// Disable the usage of real-time dynamic content.
	static bool
}

// Log implements the Formatter.Log interface.
func (consoleFormatter *consoleFormatterImpl) Log(content LogContent, level loggers.LogLevel) {
	// If dynamic content is allowed, and the content implements the LogDynamicContent interface,
	// use the dynamic rendering method.
	if !consoleFormatter.static {
		dynamic, isDynamic := content.(LogDynamicContent)
		if isDynamic {
			dynamic.RenderConsoleDynamic(func(msg string) { consoleFormatter.logger.Log(level, msg) })
			return
		}
	}

	consoleFormatter.logger.Log(level, content.RenderConsole())
}

// NewConsoleFormatter creates a new Formatter that will dispatch messages to a logger.Console.
//
// If static is set to true, the formatter will ignore dynamic content.
func NewConsoleFormatter(consoleLogger loggers.Console, static bool) Formatter {
	return &consoleFormatterImpl{logger: consoleLogger, static: static}
}

// Default implementation of the Formatter interface, for logger.JSON instances.
type jsonFormatterImpl struct {
	logger loggers.JSON
}

// Log implements the Formatter.Log interface.
func (jsonFormatter *jsonFormatterImpl) Log(content LogContent, level loggers.LogLevel) {
	jsonFormatter.logger.Log(level, content.RenderJSON())
}

// NewJSONFormatter creates a new Formatter that will dispatch messages to a logger.JSON.
func NewJSONFormatter(jsonLogger loggers.JSON) Formatter {
	return &jsonFormatterImpl{logger: jsonLogger}
}

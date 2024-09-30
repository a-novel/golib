package formatters

// LogSplit allows to set specific logs format for console and JSON.
type LogSplit interface {
	LogContent

	// SetConsoleMessage sets a message that will only be rendered if the logger implements loggers.Console.
	SetConsoleMessage(content string) LogSplit
	// SetJSONMessage sets a message that will only be rendered if the logger implements loggers.JSON.
	SetJSONMessage(content interface{}) LogSplit
	// SetConsoleContent sets a LogContent that will only be rendered if the logger implements loggers.Console.
	SetConsoleContent(content LogContent) LogSplit
	// SetJSONContent sets a LogContent that will only be rendered if the logger implements loggers.JSON.
	SetJSONContent(content LogContent) LogSplit
	// SetConsoleRenderer sets a render function that will only be triggered if the logger implements loggers.Console.
	SetConsoleRenderer(content func() string) LogSplit
	// SetJSONRenderer sets a render function that will only be triggered if the logger implements loggers.JSON.
	SetJSONRenderer(content func() interface{}) LogSplit
}

// Default implementation of the LogSplit interface.
type logSplitImpl struct {
	consoleMessage func() string
	jsonMessage    func() interface{}
}

// RenderConsole implements LogContent.RenderConsole interface.
func (logSplit *logSplitImpl) RenderConsole() string {
	if logSplit.consoleMessage == nil {
		return ""
	}

	return logSplit.consoleMessage()
}

// RenderJSON implements LogContent.RenderJSON interface.
func (logSplit *logSplitImpl) RenderJSON() interface{} {
	if logSplit.jsonMessage == nil {
		return nil
	}

	return logSplit.jsonMessage()
}

// SetConsoleMessage implements LogSplit.SetConsoleMessage interface.
func (logSplit *logSplitImpl) SetConsoleMessage(content string) LogSplit {
	logSplit.consoleMessage = func() string { return content + "\n" }
	return logSplit
}

// SetJSONMessage implements LogSplit.SetJSONMessage interface.
func (logSplit *logSplitImpl) SetJSONMessage(content interface{}) LogSplit {
	logSplit.jsonMessage = func() interface{} { return content }
	return logSplit
}

// SetConsoleContent implements LogSplit.SetConsoleContent interface.
func (logSplit *logSplitImpl) SetConsoleContent(content LogContent) LogSplit {
	logSplit.consoleMessage = content.RenderConsole
	return logSplit
}

// SetConsoleRenderer implements LogSplit.SetConsoleRenderer interface.
func (logSplit *logSplitImpl) SetConsoleRenderer(content func() string) LogSplit {
	logSplit.consoleMessage = content
	return logSplit
}

// SetJSONRenderer implements LogSplit.SetJSONRenderer interface.
func (logSplit *logSplitImpl) SetJSONRenderer(content func() interface{}) LogSplit {
	logSplit.jsonMessage = content
	return logSplit
}

// SetJSONContent implements LogSplit.SetJSONContent interface.
func (logSplit *logSplitImpl) SetJSONContent(content LogContent) LogSplit {
	logSplit.jsonMessage = content.RenderJSON
	return logSplit
}

// NewSplit creates a new LogSplit instance.
func NewSplit() LogSplit {
	return &logSplitImpl{}
}

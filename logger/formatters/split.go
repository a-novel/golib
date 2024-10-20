package formatters

// LogSplit allows to set specific logs format for console and JSON.
type LogSplit interface {
	LogContent
	SetConsoleMessage(content string) LogSplit
	SetJSONMessage(content interface{}) LogSplit
	SetConsoleContent(content LogContent) LogSplit
	SetJSONContent(content interface{}) LogSplit
}

// Default implementation of the LogSplit interface.
type logSplitImpl struct {
	consoleMessage string
	jsonMessage    interface{}
}

// RenderConsole implements LogContent.RenderConsole interface.
func (l *logSplitImpl) RenderConsole() string {
	return l.consoleMessage
}

// RenderJSON implements LogContent.RenderJSON interface.
func (l *logSplitImpl) RenderJSON() interface{} {
	return l.jsonMessage
}

// SetConsoleMessage implements LogSplit.SetConsoleMessage interface.
func (l *logSplitImpl) SetConsoleMessage(content string) LogSplit {
	l.consoleMessage = content
	return l
}

// SetJSONMessage implements LogSplit.SetJSONMessage interface.
func (l *logSplitImpl) SetJSONMessage(content interface{}) LogSplit {
	l.jsonMessage = content
	return l
}

// SetConsoleContent implements LogSplit.SetConsoleContent interface.
func (l *logSplitImpl) SetConsoleContent(content LogContent) LogSplit {
	l.consoleMessage = content.RenderConsole()
	return l
}

// SetJSONContent implements LogSplit.SetJSONContent interface.
func (l *logSplitImpl) SetJSONContent(content interface{}) LogSplit {
	l.jsonMessage = content
	return l
}

// NewSplit creates a new LogSplit instance.
func NewSplit() LogSplit {
	return &logSplitImpl{}
}

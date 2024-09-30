package formatters

import "github.com/charmbracelet/lipgloss"

// LogBase is a simple implementation of LogContent. It supports colored output using lipgloss.Style, and
// can contain a child LogContent.
type LogBase interface {
	LogContent
	// SetStyle sets the style of the message. The style is only rendered with logger.ConsoleLogger.
	SetStyle(style lipgloss.Style) LogBase
	// SetChild sets a child LogContent to the message.
	SetChild(child LogContent) LogBase
}

// Default implementation of LogBase.
type logBaseImpl struct {
	content string
	style   *lipgloss.Style
	child   LogContent
}

// SetStyle implements LogBase.SetStyle interface.
func (logBase *logBaseImpl) SetStyle(style lipgloss.Style) LogBase {
	logBase.style = &style
	return logBase
}

// SetChild implements LogBase.SetChild interface.
func (logBase *logBaseImpl) SetChild(child LogContent) LogBase {
	logBase.child = child
	return logBase
}

// RenderConsole implements LogContent.RenderConsole interface.
func (logBase *logBaseImpl) RenderConsole() string {
	// Set the child string. If no string, it will remain empty, and add nothing to the log.
	child := ""
	if logBase.child != nil {
		// Render the child content, and append it to the base content.
		child = logBase.child.RenderConsole()
	}

	// Render the content with the style, if any.
	if logBase.style != nil {
		return logBase.style.Render(logBase.content) + "\n" + child
	}

	return logBase.content + "\n" + child
}

// RenderJSON implements LogContent.RenderJSON interface.
func (logBase *logBaseImpl) RenderJSON() interface{} {
	output := map[string]interface{}{"message": logBase.content}

	// Add the child content to the output, if any.
	if logBase.child != nil {
		output["child"] = logBase.child.RenderJSON()
	}

	return output
}

// NewBase creates a new LogBase instance.
func NewBase(content string) LogBase {
	return &logBaseImpl{content: content}
}

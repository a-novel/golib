package formatters

import "github.com/charmbracelet/lipgloss"

// LogBase is a simple implementation of LogContent. It supports colored output using lipgloss.Style, and
// can contain a child LogContent.
type LogBase interface {
	LogContent
	// SetStyle sets the style of the message. The style is only rendered with logger.ConsoleLogger.
	SetStyle(style *lipgloss.Style) LogBase
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
func (l *logBaseImpl) SetStyle(style *lipgloss.Style) LogBase {
	l.style = style
	return l
}

// SetChild implements LogBase.SetChild interface.
func (l *logBaseImpl) SetChild(child LogContent) LogBase {
	l.child = child
	return l
}

// RenderConsole implements LogContent.RenderConsole interface.
func (l *logBaseImpl) RenderConsole() string {
	// Set the child string. If no string, it will remain empty, and add nothing to the log.
	child := ""
	if l.child != nil {
		// Render the child content, and append it to the base content.
		child = l.child.RenderConsole()
	}

	// Render the content with the style, if any.
	if l.style != nil {
		return l.style.Render(l.content) + "\n" + child
	}

	return l.content + "\n" + child
}

// RenderJSON implements LogContent.RenderJSON interface.
func (l *logBaseImpl) RenderJSON() interface{} {
	output := map[string]interface{}{"message": l.content}

	// Add the child content to the output, if any.
	if l.child != nil {
		output["child"] = l.child.RenderJSON()
	}

	return output
}

// NewBase creates a new LogBase instance.
func NewBase(content string) LogBase {
	return &logBaseImpl{content: content}
}

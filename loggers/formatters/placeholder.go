package formatters

import "github.com/charmbracelet/lipgloss"

// LogPlaceholder renders a special message for indicating a no-op.
type LogPlaceholder interface {
	LogContent
}

// Default implementation of LogPlaceholder.
type logPlaceholderImpl struct {
	content string
}

// RenderConsole implements LogContent.RenderConsole interface.
func (logPlaceholder *logPlaceholderImpl) RenderConsole() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8000")).Render("⚠  ") +
		lipgloss.NewStyle().Faint(true).Render(logPlaceholder.content) + "\n"
}

// RenderJSON implements LogContent.RenderJSON interface.
func (logPlaceholder *logPlaceholderImpl) RenderJSON() interface{} {
	return map[string]interface{}{"message": logPlaceholder.content}
}

// NewPlaceholder creates a new LogPlaceholder instance.
func NewPlaceholder(content string) LogPlaceholder {
	return &logPlaceholderImpl{content: content}
}

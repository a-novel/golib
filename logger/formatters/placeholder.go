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
func (l *logPlaceholderImpl) RenderConsole() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8000")).Render("⚠  ") +
		lipgloss.NewStyle().Faint(true).Render(l.content) + "\n"
}

// RenderJSON implements LogContent.RenderJSON interface.
func (l *logPlaceholderImpl) RenderJSON() interface{} {
	return map[string]interface{}{"message": l.content}
}

// NewPlaceholder creates a new LogPlaceholder instance.
func NewPlaceholder(content string) LogPlaceholder {
	return &logPlaceholderImpl{content: content}
}

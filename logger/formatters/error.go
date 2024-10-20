package formatters

import (
	"github.com/charmbracelet/lipgloss"
)

// LogError renders an error as a message.
type LogError interface {
	LogContent
}

// Default implementation of LogError.
type logErrorImpl struct {
	err     error
	message string
}

// RenderConsole returns a string representation of the log content, that can be printed in a logger.ConsoleLogger.
func (l *logErrorImpl) RenderConsole() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3232")).Render(l.message + ": " + l.err.Error())
}

// RenderJSON returns a serializable representation of the log content, that can be printed in a logger.JSONLogger.
func (l *logErrorImpl) RenderJSON() interface{} {
	return map[string]interface{}{
		"error":   l.err.Error(),
		"message": l.message,
	}
}

// NewLogError creates a new LogError instance.
func NewLogError(err error, message string) LogError {
	return &logErrorImpl{
		err:     err,
		message: message,
	}
}

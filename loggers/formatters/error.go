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
func (logError *logErrorImpl) RenderConsole() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF3232")).
		Render(logError.message+": "+logError.err.Error()) + "\n"
}

// RenderJSON returns a serializable representation of the log content, that can be printed in a logger.JSONLogger.
func (logError *logErrorImpl) RenderJSON() interface{} {
	return map[string]interface{}{
		"error":   logError.err.Error(),
		"message": logError.message,
	}
}

// NewError creates a new LogError instance.
func NewError(err error, message string) LogError {
	return &logErrorImpl{
		err:     err,
		message: message,
	}
}

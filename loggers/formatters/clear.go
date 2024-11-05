package formatters

import (
	"github.com/charmbracelet/x/ansi"
)

// LogClear is a utility LogContent that clears the terminal screen on logger.ConsoleLogger instances.
// It does nothing on logger.JSONLogger instances.
type LogClear interface {
	LogContent
}

// Default implementation of LogClear.
type logClearImpl struct{}

// RenderConsole implements LogContent.RenderConsole interface.
func (logClear *logClearImpl) RenderConsole() string {
	return ansi.EraseDisplay(2) + ansi.HomeCursorPosition
}

// RenderJSON implements LogContent.RenderJSON interface.
func (logClear *logClearImpl) RenderJSON() interface{} {
	// Since there is no point in clearing the screen on a JSON logger, simply return nil.
	// logger.JSONLogger should ignore nil inputs, and treat this as a no-op.
	return nil
}

// NewClear creates a new LogClear instance.
func NewClear() LogClear {
	return &logClearImpl{}
}

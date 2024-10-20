package logger

// Implements Console using the io.writer interface.
type dummyLogger struct{}

// Log implements the Console.Log interface.
func (l *dummyLogger) Log(_ LogLevel, _ string) {}

// NewDummy creates a new Console that does nothing.
func NewDummy() Console {
	return &dummyLogger{}
}

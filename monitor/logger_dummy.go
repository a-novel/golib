package anovelmonitor

import (
	"context"
	"github.com/gin-gonic/gin"
)

// Holds the partialLogger implementation.
type dummyLogger struct{}

func (l *dummyLogger) Fatal(_ string, _ ...interface{}) {}

func (l *dummyLogger) Error(_ string, _ ...interface{}) {}

func (l *dummyLogger) Warn(_ string, _ ...interface{}) {}

func (l *dummyLogger) Info(_ string, _ ...interface{}) {}

func (l *dummyLogger) Write(_ []byte) (int, error) { return 0, nil }

// NewDummyLogger returns a no-op logger, that can be used as a placeholder for tests.
func NewDummyLogger() Logger {
	return newLoggerFromImpl(&dummyLogger{})
}

// Holds the GinLogger implementation.
type dummyGinLogger struct {
	Logger
}

func (l *dummyGinLogger) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// NewDummyGinLogger returns a no-op GinLogger, that can be used as a placeholder for tests.
func NewDummyGinLogger() GinLogger {
	return &dummyGinLogger{NewDummyLogger()}
}

// Holds the GRPCLogger implementation.
type dummyGRPCLogger struct {
	Logger
}

func (l *dummyGRPCLogger) Report(_ context.Context, _ string, _ error) {}

// NewDummyGRPCLogger returns a no-op GRPCLogger, that can be used as a placeholder for tests.
func NewDummyGRPCLogger() GRPCLogger {
	return &dummyGRPCLogger{NewDummyLogger()}
}

package anovelmonitor

import (
	"context"
	"github.com/gin-gonic/gin"
	"io"
)

// Logger represents a format agnostic logger to manage monitoring across environments.
type Logger interface {
	// FatalE quickly prints error, with an optional descriptive message.
	FatalE(err error, msg string)
	// ErrorE quickly prints error, with an optional descriptive message.
	ErrorE(err error, msg string)

	// Fatal logs to stderr and exits the program.
	Fatal(msg string, args ...interface{})
	// Error prints to stderr.
	Error(msg string, args ...interface{})
	// Warn prints to a dedicated channel if available, otherwise it should redirect to the standard output.
	// If q dedicated channel is missing, the logger can add some formatting to distinguish it from the standard output.
	Warn(msg string, args ...interface{})
	// Info prints to stdout.
	Info(msg string, args ...interface{})

	// Writer allows the logger to be used in place of the default log package. Useful when context is limited, such
	// as during startup.
	io.Writer
}

// GinLogger extends the capacities of Logger, with a dedicated middleware for Gin. Errors are automatically retrieved
// from the context, and log level is inferred from the status code.
type GinLogger interface {
	// Logger extends the base interface, so GinLogger can still be used as a classic logger, without allocating extra
	// variables.
	Logger

	// Middleware returns a Gin middleware that automatically logs requests and responses.
	// This middleware must be located at the beginning of the middleware chain.
	Middleware() gin.HandlerFunc
}

// GRPCLogger extends the capacities of Logger, with a dedicated capture method for a GRPC service. The log level
// is inferred from the GRPC code contained in the error (Logger.Info by default).
type GRPCLogger interface {
	// Logger extends the base interface, so GRPCLogger can still be used as a classic logger, without allocating extra
	// variables.
	Logger

	// Report is called on the output of a GRPC service, along with the service name. It automatically creates a
	// monitoring log for the current trigger.
	Report(ctx context.Context, service string, err error)
}

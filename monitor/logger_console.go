package anovelmonitor

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
	"time"
)

// Holds the partialLogger implementation.
type consoleLogger struct{}

func (l *consoleLogger) Fatal(msg string, args ...interface{}) {
	// Use magenta color to distinguish fatal crashes from regular errors.
	colorizer := color.New(color.FgMagenta).SprintFunc()
	log.Fatalln(colorizer(fmt.Sprintf(msg, args...)))
}

func (l *consoleLogger) Error(msg string, args ...interface{}) {
	colorizer := color.New(color.FgRed).SprintFunc()
	log.Println(colorizer(fmt.Sprintf(msg, args...)))
}

func (l *consoleLogger) Warn(msg string, args ...interface{}) {
	colorizer := color.New(color.FgYellow).SprintFunc()
	log.Println(colorizer(fmt.Sprintf(msg, args...)))
}

func (l *consoleLogger) Info(msg string, args ...interface{}) {
	log.Println(fmt.Sprintf(msg, args...))
}

func (l *consoleLogger) Write(p []byte) (n int, err error) {
	log.Println(string(p))
	return len(p), nil
}

// NewConsoleLogger creates a new Logger that prints to the console.
// Every output is redirected to stdout. This logger is not suitable for production monitoring, but rather tailored
// for easy local debugging.
// Log levels are differentiated using colors, rather than separate channels.
func NewConsoleLogger() Logger {
	return newLoggerFromImpl(&consoleLogger{})
}

// Holds the GinLogger implementation.
type consoleGinLogger struct {
	Logger
}

func (l *consoleGinLogger) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Compute timestamps to measure the time taken by the request.
		start := time.Now()
		c.Next()
		end := time.Now()

		// Set up color depending on status code. Add a cute prefix for even better readability.
		colorizer := color.New(color.FgBlue).SprintFunc()
		prefix := "✓"
		if c.Writer.Status() > 499 {
			colorizer = color.New(color.FgRed).SprintFunc()
			prefix = "✗"
		} else if c.Writer.Status() > 399 || len(c.Errors) > 0 {
			colorizer = color.New(color.FgYellow).SprintFunc()
			prefix = "⟁"
		}

		// Generate a human-readable message, with important information at glance.
		message := strings.Join([]string{
			"-",
			colorizer(color.New(color.Bold).Sprintf("%s %v", prefix, c.Writer.Status())),
			colorizer(fmt.Sprintf("[%s %s]", c.Request.Method, c.FullPath())),
			color.New(color.Faint).Sprint(fmt.Sprintf("(processed in %s)", end.Sub(start))),
		}, " ")

		// Feed the parent logger with relevant information.
		// Since the console logger does not make distingo between stderr and stdout (go log packages redirects
		// everything to stdout), we use the Logger.Info channel for every log. Colorization is used for level
		// differentiation.
		l.Info(message)
		for _, err := range c.Errors {
			l.ErrorE(err, "")
		}
	}
}

// NewConsoleGinLogger creates a new GinLogger implementation of NewConsoleLogger.
func NewConsoleGinLogger() GinLogger {
	return &consoleGinLogger{NewConsoleLogger()}
}

// Holds the GRPCLogger implementation.
type consoleGRPCLogger struct {
	Logger
}

func (l *consoleGRPCLogger) Report(_ context.Context, service string, err error) {
	// Set up color depending on status code. Add a cute prefix for even better readability.
	colorizer := color.New(color.FgBlue).SprintFunc()
	prefix := "✓"
	code := codes.OK
	if err != nil {
		code = status.Code(err)

		if code == codes.Unavailable || code == codes.Canceled || code == codes.Unimplemented {
			// Reserve special (yellow) treatment for codes that likely indicate a service / implementation issues.
			// GRPC does not have clear distinction between server-side and client-side errors (such as HTTP).
			colorizer = color.New(color.FgYellow).SprintFunc()
			prefix = "⟁"
		} else {
			// Regular codes are treated as standard errors.
			colorizer = color.New(color.FgRed).SprintFunc()
			prefix = "✗"
		}
	}

	// Generate a human-readable message, with important information at glance.
	message := strings.Join([]string{
		"-",
		colorizer(color.New(color.Bold).Sprintf("%s %s", prefix, code)),
		colorizer(fmt.Sprintf("[%s]", service)),
	}, " ")

	// Feed the parent logger with relevant information.
	// Since the console logger does not make distingo between stderr and stdout (go log packages redirects
	// everything to stdout), we use the Logger.Info channel for every log. Colorization is used for level
	// differentiation.
	l.Info(message)
	if err != nil {
		l.ErrorE(err, "")
	}
}

// NewConsoleGRPCLogger creates a new GRPCLogger implementation of NewConsoleLogger.
func NewConsoleGRPCLogger() GRPCLogger {
	return &consoleGRPCLogger{NewConsoleLogger()}
}

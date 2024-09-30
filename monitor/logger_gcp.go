package anovelmonitor

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

// Holds the Logger implementation.
type gcpLogger struct {
	// We use zerolog to optimize performances of structured JSON outputs.
	logger zerolog.Logger
	// The cloud run project ID, used to filter logs in the GCP console.
	projectID string
}

func (l *gcpLogger) FatalE(err error, msg string) {
	// Use the dedicated zerolog error handler, rather than relying on the agnostic partialLogger one.
	l.logger.Fatal().Err(err).Msg(msg)
}

func (l *gcpLogger) ErrorE(err error, msg string) {
	// Use the dedicated zerolog error handler, rather than relying on the agnostic partialLogger one.
	l.logger.Error().Err(err).Msg(msg)
}

func (l *gcpLogger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatal().Msgf(msg, args...)
}

func (l *gcpLogger) Error(msg string, args ...interface{}) {
	l.logger.Error().Msgf(msg, args...)
}

func (l *gcpLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warn().Msgf(msg, args...)
}

func (l *gcpLogger) Info(msg string, args ...interface{}) {
	l.logger.Info().Msgf(msg, args...)
}

func (l *gcpLogger) Write(p []byte) (n int, err error) {
	l.logger.Info().Msg(string(p))
	return len(p), nil
}

// Since this logger does not use the generic partialLogger implementation, we use a private constructor, shared by
// every gcp imnplementation.
func newGCPLogger(logger zerolog.Logger, projectID string) *gcpLogger {
	return &gcpLogger{
		logger:    logger,
		projectID: projectID,
	}
}

// NewGCPLogger creates a new Logger tailored for GCP environments. It returns structured logs, ideal for analysis and
// filtering from the GCP console.
func NewGCPLogger(logger zerolog.Logger, projectID string) Logger {
	return newGCPLogger(logger, projectID)
}

// Holds the GinLogger implementation.
type gcpGinLogger struct {
	gcpLogger
}

func (l *gcpGinLogger) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Compute timestamps to measure the time taken by the request.
		start := time.Now()
		c.Next()
		end := time.Now()

		// Infer log level from the status code.
		logLevel := zerolog.TraceLevel
		severity := "INFO" // For GCP, internally.
		if c.Writer.Status() > 499 {
			logLevel = zerolog.ErrorLevel
			severity = "ERROR"
		} else if c.Writer.Status() > 399 || len(c.Errors) > 0 {
			logLevel = zerolog.WarnLevel
			severity = "WARNING"
		}

		// Retrieve query, for better analysis.
		parsedQuery := zerolog.Dict()
		for k, v := range c.Request.URL.Query() {
			parsedQuery.Strs(k, v)
		}

		// Allow logs to be grouped in log explorer.
		// https://cloud.google.com/run/docs/logging#run_manual_logging-go
		var trace string
		if l.projectID != "" {
			traceHeader := c.GetHeader("X-Cloud-Trace-Context")
			traceParts := strings.Split(traceHeader, "/")
			if len(traceParts) > 0 && len(traceParts[0]) > 0 {
				trace = fmt.Sprintf("projects/%s/traces/%s", l.projectID, traceParts[0])
			}
		}

		ll := l.logger.WithLevel(logLevel).
			Dict(
				"httpRequest", zerolog.Dict().
					Str("requestMethod", c.Request.Method).
					Str("requestUrl", c.FullPath()).
					Int("status", c.Writer.Status()).
					Str("userAgent", c.Request.UserAgent()).
					Str("remoteIp", c.ClientIP()).
					Str("protocol", c.Request.Proto).
					Str("latency", end.Sub(start).String()),
			).
			Time("start", start).
			Str("ip", c.ClientIP()).
			Str("contentType", c.ContentType()).
			Strs("errors", c.Errors.Errors()).
			Dict("query", parsedQuery).
			Str("severity", severity)

		// Add trace to the log, if available.
		if len(trace) > 0 {
			ll = ll.Str("logging.googleapis.com/trace", trace)
		}

		// Actually send message with zerolog.
		ll.Msg(c.Request.URL.String())

		// Retrieve Sentry from the context, and if available, forward errors to it.
		hub := sentrygin.GetHubFromContext(c)
		if hub != nil {
			hub.Scope().SetRequest(c.Request)
			for _, err := range c.Errors {
				hub.CaptureException(err)
			}
		}
	}
}

// NewGCPGinLogger creates a new GinLogger implementation of NewGCPLogger. It also features automated Sentry support.
// If a valid Sentry instance is configured, errors will be forwarded to it.
func NewGCPGinLogger(logger zerolog.Logger, projectID string) GinLogger {
	return &gcpGinLogger{*newGCPLogger(logger, projectID)}
}

// Holds the GRPCLogger implementation.
type gcpGRPCLogger struct {
	gcpLogger
}

func (l *gcpGRPCLogger) Report(ctx context.Context, service string, err error) {
	// Infer log level from the status code.
	logLevel := zerolog.TraceLevel
	severity := "INFO" // For GCP.
	code := codes.OK

	if err != nil {
		code = status.Code(err)

		if code == codes.Unavailable || code == codes.Canceled || code == codes.Unimplemented {
			// Reserve special (yellow) treatment for codes that likely indicate a service / implementation issues.
			// GRPC does not have clear distinction between server-side and client-side errors (such as HTTP).
			logLevel = zerolog.WarnLevel
			severity = "WARNING"
			code = status.Code(err)
		} else {
			// Regular codes are treated as standard errors.
			logLevel = zerolog.ErrorLevel
			severity = "ERROR"
			code = status.Code(err)
		}
	}

	// TODO: check if we can add trace to GRPC requests.
	// TODO: improve formatting of GRPC messages.
	ll := l.logger.WithLevel(logLevel).
		Dict(
			"grpcRequest", zerolog.Dict().
				Str("service", service).
				Uint32("code", uint32(code)),
		).
		Err(err).
		Str("severity", severity)

	// Actually send message with zerolog.
	ll.Msg(fmt.Sprintf("GRPC %s [status %s]", service, code))

	// Retrieve Sentry from the context, and if available, forward errors to it.
	hub := sentry.GetHubFromContext(ctx)
	if hub != nil && err != nil {
		hub.Scope()
		hub.CaptureException(err)
	}
}

// NewGCPGRPCLogger creates a new GRPCLogger implementation of NewGCPLogger. It also features automated Sentry support.
// If a valid Sentry instance is configured, errors will be forwarded to it.
func NewGCPGRPCLogger(logger zerolog.Logger, projectID string) GRPCLogger {
	return &gcpGRPCLogger{*newGCPLogger(logger, projectID)}
}

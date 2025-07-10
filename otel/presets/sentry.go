package otelpresets

import (
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type SentryOtelConfig struct {
	DSN          string        `json:"dsn"          yaml:"dsn"`
	ServerName   string        `json:"serverName"   yaml:"serverName"`
	Release      string        `json:"release"      yaml:"release"`
	Environment  string        `json:"environment"  yaml:"environment"`
	FlushTimeout time.Duration `json:"flushTimeout" yaml:"flushTimeout"`
	Debug        bool          `json:"debug"        yaml:"debug"`
}

func (config *SentryOtelConfig) Init() error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:              config.DSN,
		EnableTracing:    true,
		EnableLogs:       true,
		TracesSampleRate: 1.0,
		Debug:            config.Debug,
		DebugWriter:      os.Stderr,
		ServerName:       config.ServerName,
		Release:          config.Release,
		Environment:      config.Environment,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint == nil || hint.Context == nil {
				return event
			}

			if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
				// Add IP Address to user information.
				event.User.IPAddress = req.RemoteAddr
			}

			return event
		},
	})
}

func (config *SentryOtelConfig) GetPropagators() (propagation.TextMapPropagator, error) {
	return sentryotel.NewSentryPropagator(), nil
}

func (config *SentryOtelConfig) GetTraceProvider() (trace.TracerProvider, error) {
	return sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor())), nil
}

func (config *SentryOtelConfig) GetLogger() (log.LoggerProvider, error) {
	// TODO: switch to Sentry native logger for production use.
	logExporter, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}

	return sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	), nil
}

func (config *SentryOtelConfig) Flush() {
	sentry.Flush(config.FlushTimeout)
}

func (config *SentryOtelConfig) HTTPHandler() func(http.Handler) http.Handler {
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	return sentryHandler.Handle
}

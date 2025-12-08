package otelpresets

import (
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/a-novel-kit/golib/otel"
)

var _ otel.Config = (*Sentry)(nil)

type Sentry struct {
	DSN          string        `json:"dsn"          yaml:"dsn"`
	ServerName   string        `json:"serverName"   yaml:"serverName"`
	Release      string        `json:"release"      yaml:"release"`
	Environment  string        `json:"environment"  yaml:"environment"`
	FlushTimeout time.Duration `json:"flushTimeout" yaml:"flushTimeout"`
	Debug        bool          `json:"debug"        yaml:"debug"`
}

func (config *Sentry) Init() error {
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

func (config *Sentry) GetPropagators() (propagation.TextMapPropagator, error) {
	return sentryotel.NewSentryPropagator(), nil
}

func (config *Sentry) GetTraceProvider() (trace.TracerProvider, error) {
	return sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor())), nil
}

func (config *Sentry) GetLogger() (log.LoggerProvider, error) {
	// TODO: switch to Sentry native logger for production use.
	logExporter, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}

	return sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	), nil
}

func (config *Sentry) Flush() {
	sentry.Flush(config.FlushTimeout)
}

func (config *Sentry) HttpHandler() func(http.Handler) http.Handler {
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	return sentryHandler.Handle
}

func (config *Sentry) RpcInterceptor() grpc.ServerOption {
	return grpc.StatsHandler(otelgrpc.NewServerHandler())
}

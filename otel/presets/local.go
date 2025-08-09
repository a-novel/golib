package otelpresets

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	libotel "github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/otel/utils"
)

// LocalOtelConfig configures OTEL to log traces & logs to stdout.
type LocalOtelConfig struct {
	PrettyPrint  bool          `json:"prettyPrint"  yaml:"prettyPrint"`
	FlushTimeout time.Duration `json:"flushTimeout" yaml:"flushTimeout"`
}

// Init just prints a banner for local dev mode.
func (config *LocalOtelConfig) Init() error {
	green := color.New(color.FgGreen).Add(color.Bold)
	_, _ = green.Println("🚀 OpenTelemetry Local Mode: All traces and logs to stdout")

	return nil
}

func (config *LocalOtelConfig) GetPropagators() (propagation.TextMapPropagator, error) {
	return propagation.TraceContext{}, nil
}

func (config *LocalOtelConfig) GetTraceProvider() (trace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(color.Output), // writes with color support
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
	)

	return tp, nil
}

func (config *LocalOtelConfig) GetLogger() (log.LoggerProvider, error) {
	logExporter, err := stdoutlog.New(
		stdoutlog.WithPrettyPrint(),
		stdoutlog.WithWriter(color.Output),
	)
	if err != nil {
		return nil, err
	}

	return sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	), nil
}

func (config *LocalOtelConfig) Flush() {
	provider := otel.GetTracerProvider()

	tp, ok := provider.(*sdktrace.TracerProvider)
	if ok {
		err := tp.Shutdown(context.Background())
		if err != nil {
			stdlog.Fatalf("Failed to shutdown tracer provider: %v\n", err)
		}
	}
}

func (config *LocalOtelConfig) HTTPHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := libotel.Tracer().Start(r.Context(), fmt.Sprintf("[%s] %s.%s", r.Method, r.Host, r.URL.Path))
			defer span.End()

			span.SetAttributes(
				attribute.String("request.method", r.Method),
				attribute.String("request.host", r.Host),
				attribute.String("request.path", r.URL.Path),
				attribute.String("request.remote_addr", r.RemoteAddr),
			)

			wrapped := &utils.CaptureHTTPResponseWriter{ResponseWriter: w}
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			status := wrapped.Status()

			span.SetAttributes(
				attribute.Int("response.status_code", status),
				attribute.String("response.status_text", http.StatusText(status)),
			)

			switch {
			case status >= http.StatusInternalServerError:
				span.RecordError(fmt.Errorf("HTTP %d: %s", status, http.StatusText(status)))
				span.SetStatus(codes.Error, http.StatusText(status))

				_, _ = fmt.Fprintf(
					color.Output,
					color.RedString("❌ %s %s %d: %s\n"),
					r.Method, r.URL.Path, status, http.StatusText(status),
				)
			case status >= http.StatusBadRequest:
				span.SetStatus(codes.Error, http.StatusText(status))

				_, _ = fmt.Fprintf(
					color.Output,
					color.YellowString("⚠️ %s %s %d: %s\n"),
					r.Method, r.URL.Path, status, http.StatusText(status),
				)
			default:
				span.SetStatus(codes.Ok, "")

				_, _ = fmt.Fprintf(
					color.Output,
					color.GreenString("✅ %s %s %d: %s\n"),
					r.Method, r.URL.Path, status, http.StatusText(status),
				)
			}
		})
	}
}

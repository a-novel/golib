package otelpresets

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"strconv"
	"time"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	libotel "github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/otel/utils"
)

type GCloudOtelConfig struct {
	ProjectID    string        `json:"projectID"    yaml:"projectID"`
	FlushTimeout time.Duration `json:"flushTimeout" yaml:"flushTimeout"`
}

func (config *GCloudOtelConfig) Init() error {
	blue := color.New(color.FgBlue).Add(color.Bold)
	_, _ = blue.Printf("☁️ OpenTelemetry GCP Mode: exporting traces to Cloud Trace (project=%s)\n", config.ProjectID)

	return nil
}

func (config *GCloudOtelConfig) GetPropagators() (propagation.TextMapPropagator, error) {
	return propagation.NewCompositeTextMapPropagator(
		// Putting the CloudTraceOneWayPropagator first means the TraceContext propagator
		// takes precedence if both the traceparent and the XCTC headers exist.
		gcppropagator.CloudTraceOneWayPropagator{},
		propagation.TraceContext{},
		propagation.Baggage{},
	), nil
}

func (config *GCloudOtelConfig) GetTraceProvider() (trace.TracerProvider, error) {
	var opts []texporter.Option
	if config.ProjectID != "" {
		opts = append(opts, texporter.WithProjectID(config.ProjectID))
	}

	exporter, err := texporter.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCP trace exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
	)

	return tp, nil
}

// GetLogger sets up a structured JSON logger that GCP can parse and link to traces.
func (config *GCloudOtelConfig) GetLogger() (log.LoggerProvider, error) {
	zerologLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	logExporter, err := stdoutlog.New(
		stdoutlog.WithWriter(zerologLogger),
	)
	if err != nil {
		return nil, err
	}

	return sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	), nil
}

// Flush shuts down tracer and logger providers.
func (config *GCloudOtelConfig) Flush() {
	provider := otel.GetTracerProvider()

	tp, ok := provider.(*sdktrace.TracerProvider)
	if ok {
		err := tp.Shutdown(context.Background())
		if err != nil {
			stdlog.Fatalf("Failed to shutdown tracer provider: %v\n", err)
		}
	}
}

// HTTPHandler is kept for interface compatibility but doesn't wrap anything in GCP mode.
func (config *GCloudOtelConfig) HTTPHandler() func(http.Handler) http.Handler {
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

			start := time.Now()

			next.ServeHTTP(wrapped, r.WithContext(ctx))

			latency := time.Since(start)
			status := wrapped.Status()

			span.SetAttributes(
				attribute.Int("response.status_code", status),
				attribute.String("response.status_text", http.StatusText(status)),
			)

			var logLevel string

			switch {
			case status >= http.StatusInternalServerError:
				span.RecordError(fmt.Errorf("HTTP %d: %s", status, http.StatusText(status)))
				span.SetStatus(codes.Error, http.StatusText(status))

				logLevel = "ERROR"
			case status >= http.StatusBadRequest:
				span.SetStatus(codes.Error, http.StatusText(status))

				logLevel = "WARNING"
			default:
				span.SetStatus(codes.Ok, "")

				logLevel = "INFO"
			}

			// Extract trace info for GCP
			spanCtx := span.SpanContext()
			traceID := spanCtx.TraceID().String()
			spanID := spanCtx.SpanID().String()
			traceSampled := spanCtx.IsSampled()

			// GCP trace field format
			traceResource := fmt.Sprintf("projects/%s/traces/%s", config.ProjectID, traceID)

			// Build structured log entry
			logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
			logger.Info().
				Str("severity", logLevel).
				Str("logging.googleapis.com/trace", traceResource).
				Str("logging.googleapis.com/spanId", spanID).
				Bool("logging.googleapis.com/trace_sampled", traceSampled).
				Dict("httpRequest", zerolog.Dict().
					Str("requestMethod", r.Method).
					Str("requestUrl", r.URL.String()).
					Int("status", status).
					Int64("requestSize", r.ContentLength).
					Str("remoteIp", r.RemoteAddr).
					Str("userAgent", r.UserAgent()).
					Str("referer", r.Referer()).
					Str("protocol", r.Proto).
					Str("latency", fmt.Sprintf("%.9fs", latency.Seconds())).
					Str("responseSize", strconv.FormatInt(wrapped.Size(), 10)),
				).
				Msgf("%s %s %d", r.Method, r.URL.Path, status)
		})
	}
}

package loggingpresets

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/codes"

	"github.com/a-novel-kit/golib/logging"
	libotel "github.com/a-novel-kit/golib/otel"
	"github.com/a-novel-kit/golib/otel/utils"
)

var _ logging.HttpConfig = (*HttpGcloud)(nil)

type HttpGcloud struct {
	ProjectId string `json:"projectID" yaml:"projectID"`

	l *slog.Logger
}

func (logger *HttpGcloud) Logger() func(http.Handler) http.Handler {
	logger.l = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := libotel.Tracer().Start(r.Context(), fmt.Sprintf("[%s] %s.%s", r.Method, r.Host, r.URL.Path))
			defer span.End()

			wrapped := &utils.CaptureHTTPResponseWriter{ResponseWriter: w}

			start := time.Now()

			next.ServeHTTP(wrapped, r.WithContext(ctx))

			latency := time.Since(start)
			status := wrapped.Status()

			var (
				gcloudLevel string
				log         func(ctx context.Context, msg string, args ...any)
			)

			switch {
			case status >= http.StatusInternalServerError:
				span.RecordError(errors.New(string(wrapped.Response())))
				span.SetStatus(codes.Error, http.StatusText(status))

				gcloudLevel = "ERROR"
				log = logger.l.ErrorContext
			case status >= http.StatusBadRequest:
				span.SetStatus(codes.Error, http.StatusText(status))

				gcloudLevel = "WARNING"
				log = logger.l.WarnContext
			default:
				span.SetStatus(codes.Ok, http.StatusText(status))

				gcloudLevel = "INFO"
				log = logger.l.InfoContext
			}

			// Extract trace info for GCP
			spanCtx := span.SpanContext()
			traceID := spanCtx.TraceID().String()
			spanID := spanCtx.SpanID().String()
			traceSampled := spanCtx.IsSampled()

			// GCP trace field format
			traceResource := fmt.Sprintf("projects/%s/traces/%s", logger.ProjectId, traceID)

			// Build structured log entry
			// https://docs.cloud.google.com/logging/docs/structured-logging
			log(
				r.Context(),
				fmt.Sprintf("%s %s %d", r.Method, r.URL.Path, status),
				slog.String("severity", gcloudLevel),
				slog.String("logging.googleapis.com/trace", traceResource),
				slog.String("logging.googleapis.com/spanId", spanID),
				slog.Bool("logging.googleapis.com/trace_sampled", traceSampled),
				slog.Group(
					"httpRequest",
					slog.String("requestMethod", r.Method),
					slog.String("requestUrl", r.URL.String()),
					slog.Int("status", status),
					slog.Int64("requestSize", r.ContentLength),
					slog.String("remoteIp", r.RemoteAddr),
					slog.String("userAgent", r.UserAgent()),
					slog.String("referer", r.Referer()),
					slog.String("protocol", r.Proto),
					slog.String("latency", fmt.Sprintf("%.9fs", latency.Seconds())),
					slog.String("responseSize", strconv.FormatInt(wrapped.Size(), 10)),
				),
			)
		})
	}
}

package otel

import (
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var AppName string

func SetAppName(name string) {
	AppName = name
}

func Tracer(options ...trace.TracerOption) trace.Tracer {
	return otel.GetTracerProvider().Tracer(AppName, options...)
}

func Logger(options ...otelslog.Option) *slog.Logger {
	return otelslog.NewLogger(AppName, options...)
}

func ReportError(span trace.Span, err error) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	return err
}

func ReportSuccess[Resp any](span trace.Span, resp Resp) Resp {
	span.SetStatus(codes.Ok, "")

	return resp
}

func ReportSuccessNoContent(span trace.Span) {
	span.SetStatus(codes.Ok, "")
}

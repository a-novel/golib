package loggingpresets

import (
	"context"
	"log/slog"
	"runtime/debug"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel/trace"
)

func logInterceptor(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func logTraceId(ctx context.Context) grpclog.Fields {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return grpclog.Fields{"traceID", span.TraceID().String()}
	}

	return nil
}

func panicInterceptor(l *slog.Logger) func(p any) error {
	return func(p any) error {
		l.Error("recovered from panic", "panic", p, "stack", debug.Stack())

		return status.Errorf(codes.Internal, "%s", p)
	}
}

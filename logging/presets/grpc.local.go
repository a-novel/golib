package loggingpresets

import (
	"context"
	"log/slog"

	"github.com/fatih/color"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/trace"

	"github.com/a-novel/golib/logging"
)

var _ logging.RpcConfig = (*GrpcLocal)(nil)

type GrpcLocal struct {
	Component string `json:"component" yaml:"component"`

	l *slog.Logger
}

func (logger *GrpcLocal) UnaryInterceptor() grpc.UnaryServerInterceptor {
	logger.init()

	return grpclog.UnaryServerInterceptor(localInterceptor(logger.l), grpclog.WithFieldsFromContext(localILogTraceId))
}

func (logger *GrpcLocal) StreamInterceptor() grpc.StreamServerInterceptor {
	logger.init()

	return grpclog.StreamServerInterceptor(localInterceptor(logger.l), grpclog.WithFieldsFromContext(localILogTraceId))
}

func (logger *GrpcLocal) init() {
	if logger.l != nil {
		return
	}

	log := slog.New(slog.NewTextHandler(color.Output, &slog.HandlerOptions{}))
	logger.l = log.With("service", "gRPC/server", "component", logger.Component)
}

func localInterceptor(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func localILogTraceId(ctx context.Context) grpclog.Fields {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return grpclog.Fields{"traceID", span.TraceID().String()}
	}

	return nil
}

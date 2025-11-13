package loggingpresets

import (
	"context"
	"log/slog"
	"os"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/trace"

	"github.com/a-novel/golib/logging"
)

var _ logging.RpcConfig = (*GrpcGcloud)(nil)

type GrpcGcloud struct {
	Component string `json:"component" yaml:"component"`

	l *slog.Logger
}

func (logger *GrpcGcloud) UnaryInterceptor() grpc.UnaryServerInterceptor {
	logger.init()

	return grpclog.UnaryServerInterceptor(gcloudInterceptor(logger.l), grpclog.WithFieldsFromContext(gcloudLogTraceId))
}

func (logger *GrpcGcloud) StreamInterceptor() grpc.StreamServerInterceptor {
	logger.init()

	return grpclog.StreamServerInterceptor(gcloudInterceptor(logger.l), grpclog.WithFieldsFromContext(gcloudLogTraceId))
}

func (logger *GrpcGcloud) init() {
	if logger.l != nil {
		return
	}

	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))
	logger.l = log.With("service", "gRPC/server", "component", logger.Component)
}

func gcloudInterceptor(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func gcloudLogTraceId(ctx context.Context) grpclog.Fields {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return grpclog.Fields{"traceID", span.TraceID().String()}
	}

	return nil
}

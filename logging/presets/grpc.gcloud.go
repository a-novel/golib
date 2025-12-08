package loggingpresets

import (
	"log/slog"
	"os"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"

	"github.com/a-novel-kit/golib/logging"
)

var _ logging.RpcConfig = (*GrpcGcloud)(nil)

type GrpcGcloud struct {
	Component string `json:"component" yaml:"component"`

	l *slog.Logger
}

func (logger *GrpcGcloud) UnaryInterceptor() grpc.UnaryServerInterceptor {
	logger.init()

	return grpclog.UnaryServerInterceptor(logInterceptor(logger.l), grpclog.WithFieldsFromContext(logTraceId))
}

func (logger *GrpcGcloud) StreamInterceptor() grpc.StreamServerInterceptor {
	logger.init()

	return grpclog.StreamServerInterceptor(logInterceptor(logger.l), grpclog.WithFieldsFromContext(logTraceId))
}

func (logger *GrpcGcloud) PanicUnaryInterceptor() grpc.UnaryServerInterceptor {
	logger.init()

	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(panicInterceptor(logger.l)))
}

func (logger *GrpcGcloud) PanicStreamInterceptor() grpc.StreamServerInterceptor {
	logger.init()

	return recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(panicInterceptor(logger.l)))
}

func (logger *GrpcGcloud) init() {
	if logger.l != nil {
		return
	}

	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))
	logger.l = log.With("service", "gRPC/server", "component", logger.Component)
}

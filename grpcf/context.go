package grpcf

import (
	"context"

	"google.golang.org/grpc"
)

func BaseContextUnaryInterceptor(
	ctxInterceptor func(context.Context) context.Context,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (any, error) {
		return handler(ctxInterceptor(ctx), req)
	}
}

type wrappedStream struct {
	grpc.ServerStream

	ctxInterceptor func(context.Context) context.Context
}

func (stream *wrappedStream) Context() context.Context {
	return stream.ctxInterceptor(stream.ServerStream.Context())
}

func BaseContextStreamInterceptor(
	ctxInterceptor func(context.Context) context.Context,
) grpc.StreamServerInterceptor {
	return func(
		srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
	) error {
		return handler(srv, &wrappedStream{ss, ctxInterceptor})
	}
}

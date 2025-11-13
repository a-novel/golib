package logging

import "google.golang.org/grpc"

type RpcConfig interface {
	UnaryInterceptor() grpc.UnaryServerInterceptor
	StreamInterceptor() grpc.StreamServerInterceptor
}

package grpcf

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	golibproto "github.com/a-novel-kit/golib/grpcf/proto/gen"
)

type echo struct {
	golibproto.UnimplementedEchoServiceServer
}

func (handler *echo) UnaryEcho(context.Context, *golibproto.UnaryEchoRequest) (*golibproto.UnaryEchoResponse, error) {
	return &golibproto.UnaryEchoResponse{
		Message: "Hello world!",
	}, nil
}

func SetEchoServers(server *grpc.Server, healthPing time.Duration) {
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(server, healthcheck)

	golibproto.RegisterEchoServiceServer(server, new(echo))

	go func() {
		// asynchronously inspect dependencies and toggle serving status as needed
		next := healthpb.HealthCheckResponse_SERVING

		for {
			healthcheck.SetServingStatus("", next)

			if next == healthpb.HealthCheckResponse_SERVING {
				next = healthpb.HealthCheckResponse_NOT_SERVING
			} else {
				next = healthpb.HealthCheckResponse_SERVING
			}

			time.Sleep(healthPing)
		}
	}()
}

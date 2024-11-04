package grpc

import (
	"context"

	"github.com/a-novel/golib/loggers/adapters"
)

type ExecService[In any, Out any] interface {
	Exec(ctx context.Context, data In) (Out, error)
}

type wrappedService[In any, Out any] struct {
	service adapters.GRPCCallback[In, Out]
}

func (s *wrappedService[In, Out]) Exec(ctx context.Context, data In) (Out, error) {
	return s.service(ctx, data)
}

func ServiceWithMetrics[In any, Out any](
	name string, service ExecService[In, Out], logger adapters.GRPC,
) ExecService[In, Out] {
	return &wrappedService[In, Out]{
		service: adapters.WrapGRPCCall(name, logger, service.Exec),
	}
}

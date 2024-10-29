package grpc

import (
	"context"
	"sync"
	"time"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// DepCheckCallbacks is a function that checks the health of a list of dependencies, and return an error for any
// failed dependency.
//
// This map must contain every dependency. Healthy dependencies must be associated to a nil error value.
type DepCheckCallbacks map[string]func() error

// DepCheckServices lists dependencies for each GRPC service. When healthcheck is triggered, a GRPC service will
// be marked as healthy only if every single one of its dependencies is also healthy.
//
// You should only list RPC services here. The generic health (empty key) is automatically handled by the system.
type DepCheckServices map[string][]string

// DepsCheck configures a health checker for a GRPC service. This healthcheck is run periodically, and will prevent
// faulty calls to unhealthy services.
type DepsCheck struct {
	Dependencies DepCheckCallbacks
	Services     DepCheckServices
}

type healthServer struct {
	healthpb.UnimplementedHealthServer
	mu        sync.RWMutex
	depsCheck *DepsCheck

	watchInterval time.Duration
}

func (server *healthServer) getServiceDeps(service string) ([]string, error) {
	server.mu.RLock()
	defer server.mu.RUnlock()

	if service == "" {
		return lo.Keys(server.depsCheck.Dependencies), nil
	}

	deps, ok := server.depsCheck.Services[service]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "unknown service %s", service)
	}

	return deps, nil
}

func (server *healthServer) getStatus(request *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	deps, err := server.getServiceDeps(request.GetService())
	if err != nil {
		return nil, err
	}

	for _, dep := range deps {
		server.mu.RLock()
		depChecker, ok := server.depsCheck.Dependencies[dep]
		server.mu.RUnlock()
		if !ok {
			return nil, status.Errorf(codes.Unknown, "unknown dependency %s", dep)
		}

		if err = depChecker(); err != nil {
			return &healthpb.HealthCheckResponse{ //nolint:nilerr
				Status: healthpb.HealthCheckResponse_NOT_SERVING,
			}, nil
		}
	}

	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (server *healthServer) Check(
	_ context.Context, request *healthpb.HealthCheckRequest,
) (*healthpb.HealthCheckResponse, error) {
	return server.getStatus(request)
}

func (server *healthServer) Watch(request *healthpb.HealthCheckRequest, stream healthpb.Health_WatchServer) error {
	// Create a tick to update health status.
	server.mu.RLock()
	ticker := time.NewTicker(server.watchInterval)
	server.mu.RUnlock()

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			statusResp, err := server.getStatus(request)
			if err != nil {
				return status.Errorf(codes.Canceled, "get service status: %s", err)
			}

			serviceStatus := statusResp.GetStatus()

			if err := stream.Send(&healthpb.HealthCheckResponse{Status: serviceStatus}); err != nil {
				return status.Errorf(codes.Canceled, "stream service status: %s", err)
			}
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "stream terminated") //nolint:wrapcheck
		}
	}
}

func NewHealthServer(depsCheck *DepsCheck, watchInterval time.Duration) healthpb.HealthServer {
	return &healthServer{
		depsCheck:     depsCheck,
		watchInterval: watchInterval,
	}
}

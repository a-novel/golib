package grpc

import (
	"fmt"
	"net"
	"time"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/loggers/formatters"
)

// DepCheckCallback is a function that checks the health of a list of dependencies, and return an error for any
// failed dependency.
//
// This map must contain every dependency. Healthy dependencies must be associated to a nil error value.
type DepCheckCallback func() map[string]error

// DepCheckServices lists dependencies for each GRPC service. When healthcheck is triggered, a GRPC service will
// be marked as healthy only if every single one of its dependencies is also healthy.
//
// You should only list RPC services here. The generic health (empty key) is automatically handled by the system.
type DepCheckServices map[string][]string

// DepsCheck configures a health checker for a GRPC service. This healthcheck is run periodically, and will prevent
// faulty calls to unhealthy services.
type DepsCheck struct {
	Dependencies DepCheckCallback
	Services     DepCheckServices
}

// StartServer starts a new GRPC server on the specified port.
//
// You must ensure to properly close the server when you are done, using the CloseGRPCServer method.
//
//	listener, server, health := deploy.StartGRPCServer(50051)
//	// Graceful shutdown.
//	defer deploy.CloseGRPCServer(listener, server)
//	// (optional, but recommended) Start healthcheck.
//	go health()
func StartServer(
	port int,
	depsCheck DepsCheck,
	formatter formatters.Formatter,
) (net.Listener, *grpc.Server, func()) {
	// Prevent accidental misconfigurations.
	if port == 0 {
		formatter.Log(formatters.NewBase("port is required"), loggers.LogLevelFatal)
	}

	// Start to listen on the provided port.
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		formatter.Log(formatters.NewError(err, "listen"), loggers.LogLevelFatal)
	}

	server := grpc.NewServer()

	// Set healthcheck.
	// https://github.com/grpc/grpc-go/blob/master/examples/features/health/server/main.go
	healthcheck := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthcheck)

	healthUpdater := func() {
		// List the statuses of every dependency. Unhealthy dependencies will report an error.
		dependencies := depsCheck.Dependencies()
		// Health of the global service. If any dependency is unhealthy, the whole service will be marked as unhealthy.
		global := true

		// Check each dependency individually. Add some logs in case of failure, so the issue can be quickly identified.
		// This loop reports errors, and sets the global service health.
		for dependency, err := range dependencies {
			if err != nil {
				formatter.Log(
					formatters.NewError(err, fmt.Sprintf("dependency check for %s failed", dependency)),
					loggers.LogLevelError,
				)
				global = false
			}
		}

		// Check services individually.
		for service, serviceDeps := range depsCheck.Services {
			// Check if any of the current service dependencies is unhealthy.
			_, hasError := lo.Find(serviceDeps, func(item string) bool {
				return dependencies[item] != nil
			})

			healthcheck.SetServingStatus(
				service,
				lo.Ternary(
					hasError,
					grpc_health_v1.HealthCheckResponse_NOT_SERVING,
					grpc_health_v1.HealthCheckResponse_SERVING,
				),
			)
		}

		healthcheck.SetServingStatus(
			"",
			lo.Ternary(
				global,
				grpc_health_v1.HealthCheckResponse_SERVING,
				grpc_health_v1.HealthCheckResponse_NOT_SERVING,
			),
		)

		// Repeat healthcheck every 5 second.
		time.Sleep(5 * time.Second)
	}

	return listener, server, healthUpdater
}

// CloseServer closes an existing GRPC server.
func CloseServer(listener net.Listener, server *grpc.Server) {
	server.GracefulStop()
	_ = listener.Close()
}

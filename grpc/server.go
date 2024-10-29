package grpc

import (
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

var ErrPortRequired = errors.New("port is required")

// StartServer starts a new GRPC server on the specified port.
//
// You must ensure to properly close the server when you are done, using the CloseGRPCServer method.
//
//	listener, server := deploy.StartGRPCServer(50051)
//	// Graceful shutdown.
//	defer deploy.CloseGRPCServer(listener, server)
func StartServer(port int) (net.Listener, *grpc.Server, error) {
	// Prevent accidental misconfigurations.
	if port == 0 {
		return nil, nil, ErrPortRequired
	}

	// Start to listen on the provided port.
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, nil, fmt.Errorf("listen: %w", err)
	}

	server := grpc.NewServer()

	return listener, server, nil
}

// CloseServer closes an existing GRPC server.
func CloseServer(listener net.Listener, server *grpc.Server) {
	server.GracefulStop()
	_ = listener.Close()
}

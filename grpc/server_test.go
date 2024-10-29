package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	testgrpc "google.golang.org/grpc/interop/grpc_testing"

	anovelgrpc "github.com/a-novel/golib/grpc"
	grpcmocks "github.com/a-novel/golib/grpc/mocks"
)

func setupServerStubServer(t *testing.T) *grpcmocks.StubServer {
	t.Helper()

	return &grpcmocks.StubServer{
		EmptyCallF: func(ctx context.Context, _ *testgrpc.Empty) (*testgrpc.Empty, error) {
			return new(testgrpc.Empty), nil
		},
	}
}

func TestServerServing(t *testing.T) {
	listener, server, err := anovelgrpc.StartServer(8080)
	require.NoError(t, err)
	defer anovelgrpc.CloseServer(listener, server)

	testgrpc.RegisterTestServiceServer(server, setupServerStubServer(t))

	go func() {
		require.NoError(t, server.Serve(listener))
	}()

	connPool := anovelgrpc.NewConnPool()
	defer connPool.Close()

	conn, err := connPool.Open("127.0.0.1", 8080, anovelgrpc.ProtocolHTTPS)
	require.NoError(t, err)

	client := testgrpc.NewTestServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.EmptyCall(ctx, new(testgrpc.Empty))
	require.NoError(t, err)
}

func TestServerNoPort(t *testing.T) {
	_, _, err := anovelgrpc.StartServer(0)
	require.ErrorIs(t, err, anovelgrpc.ErrPortRequired)
}

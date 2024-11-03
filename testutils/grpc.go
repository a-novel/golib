package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// RequireGRPCCodesEqual checks if an error is nil or if it is a gRPC error with the expected code.
//
// If the code is codes.OK, it will instead check if the error is nil.
func RequireGRPCCodesEqual(t *testing.T, err error, code codes.Code) {
	t.Helper()

	if code != codes.OK {
		require.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, code, st.Code(), err.Error())
	} else {
		require.NoError(t, err)
	}
}

// WaitConn expects a conn to enter ready state. It returns when conn is available, otherwise it fails the test.
func WaitConn(t *testing.T, conn *grpc.ClientConn) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	healthClient := healthpb.NewHealthClient(conn)
	require.Eventually(
		t,
		func() bool {
			res, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{})
			return err == nil && res.GetStatus() == healthpb.HealthCheckResponse_SERVING
		},
		10*time.Second, 100*time.Millisecond,
	)
}

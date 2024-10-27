package testutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
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

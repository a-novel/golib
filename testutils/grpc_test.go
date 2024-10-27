package testutils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel/golib/testutils"
)

func TestRequireGRPCCodesMatch(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			testutils.RequireGRPCCodesEqual(t, status.Error(codes.Internal, "internal"), codes.Internal)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
		},
	})
}

func TestRequireGRPCCodesMismatch(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			testutils.RequireGRPCCodesEqual(t, status.Error(codes.Internal, "internal"), codes.NotFound)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.False(t, res.Success)
		},
	})
}

func TestRequireGRPCCodesMatchOK(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			testutils.RequireGRPCCodesEqual(t, nil, codes.OK)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
		},
	})
}

func TestRequireGRPCCodesMismatchOK(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			testutils.RequireGRPCCodesEqual(t, fmt.Errorf(""), codes.OK)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.False(t, res.Success)
		},
	})
}

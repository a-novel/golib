package testutils_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/testutils"
)

func TestRunCMDSuccess(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.True(t, true)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
		},
	})
}

func TestRunCMDFailure(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.True(t, false)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.False(t, res.Success)
		},
	})
}

func TestRunCMDCaptureSTD(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			_, _ = fmt.Fprintln(os.Stdout, "foo")
			_, _ = fmt.Fprintln(os.Stderr, "bar")
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
			require.Equal(t, "foo\n", res.STDOut)
			require.Equal(t, "bar\n", res.STDErr)
		},
	})
}

func TestRunCMDWithEnv(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, os.Getenv("FOO"), "bar")
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
		},
		Env: []string{"FOO=bar"},
	})
}

func TestRunCMDIgnoreReservedEnv(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, os.Getenv("JUST_CHECKING"), "bruh")
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success, res.STDErr)
		},
		Env: []string{"JUST_CHECKING=brother"},
	})
}

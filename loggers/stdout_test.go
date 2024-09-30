package loggers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/testutils"
)

func TestSTDOut(t *testing.T) {
	testCases := []struct {
		name string

		level loggers.LogLevel
		msg   string

		expectSTDOut string
		expectSTDErr string
	}{
		{
			name:  "Info",
			level: loggers.LogLevelInfo,
			msg:   "foo",

			expectSTDOut: "foo",
		},
		{
			name:  "Warning",
			level: loggers.LogLevelWarning,
			msg:   "foo",

			expectSTDOut: "foo",
		},
		{
			name:  "Error",
			level: loggers.LogLevelError,
			msg:   "foo",

			expectSTDErr: "foo",
		},
		{
			name:  "NoMessage",
			level: loggers.LogLevelInfo,
		},
		{
			name:  "UnknownLevel",
			level: loggers.LogLevel("unknown"),
			msg:   "foo",

			expectSTDOut: "foo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			capturer, err := captureSTD(t)
			require.NoError(t, err)

			logger := loggers.NewSTDOut()
			logger.Log(tc.level, tc.msg)

			out := capturer()

			require.Equal(t, tc.expectSTDOut, out.stdout)
			require.Equal(t, tc.expectSTDErr, out.stderr)
		})
	}
}

func TestSTDOutPanic(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			logger := loggers.NewSTDOut()
			logger.Log(loggers.LogLevelFatal, "oopsie doopsie")
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.False(t, res.Success)
			require.Equal(t, res.STDErr, "oopsie doopsie")
		},
	})
}

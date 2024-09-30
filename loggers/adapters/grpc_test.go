package adapters_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/loggers/adapters"
	adaptersmocks "github.com/a-novel/golib/loggers/adapters/mocks"
	"github.com/a-novel/golib/loggers/formatters"
	formattersmocks "github.com/a-novel/golib/loggers/formatters/mocks"
)

func TestGRPC(t *testing.T) {
	testCases := []struct {
		name string

		service string
		err     error

		expectLogLevel loggers.LogLevel
		expectConsole  *regexp.Regexp
		expectJSON     interface{}
	}{
		{
			name:    "SimpleRequest",
			service: "foo",

			expectLogLevel: loggers.LogLevelInfo,
			expectConsole:  regexp.MustCompile("^✓ OK \\[foo]\n\n$"),
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.OK,
					"service": "foo",
				},
				"severity": "INFO",
			},
		},
		{
			name:    "NonRPCError",
			service: "foo",
			err:     errors.New("oopsie"),

			expectLogLevel: loggers.LogLevelError,
			expectConsole:  regexp.MustCompile("^✗ Unknown \\[foo]\n\n  - oopsie\n\n$"),
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.Unknown,
					"service": "foo",
				},
				"severity": "ERROR",
				"error":    "oopsie",
			},
		},
		{
			name:    "RPCError",
			service: "foo",
			err:     status.Errorf(codes.NotFound, "oopsie"),

			expectLogLevel: loggers.LogLevelError,
			expectConsole:  regexp.MustCompile("^✗ NotFound \\[foo]\n\n  - rpc error: code = NotFound desc = oopsie\n\n$"),
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.NotFound,
					"service": "foo",
				},
				"severity": "ERROR",
				"error":    "rpc error: code = NotFound desc = oopsie",
			},
		},

		{
			name:    "Unavailable",
			service: "foo",
			err:     status.Errorf(codes.Unavailable, "oopsie"),

			expectLogLevel: loggers.LogLevelWarning,
			expectConsole:  regexp.MustCompile("^⟁ Unavailable \\[foo]\n\n  - rpc error: code = Unavailable desc = oopsie\n\n$"),
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.Unavailable,
					"service": "foo",
				},
				"severity": "WARNING",
				"error":    "rpc error: code = Unavailable desc = oopsie",
			},
		},
		{
			name:    "Canceled",
			service: "foo",
			err:     status.Errorf(codes.Canceled, "oopsie"),

			expectLogLevel: loggers.LogLevelWarning,
			expectConsole:  regexp.MustCompile("^⟁ Canceled \\[foo]\n\n  - rpc error: code = Canceled desc = oopsie\n\n$"),
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.Canceled,
					"service": "foo",
				},
				"severity": "WARNING",
				"error":    "rpc error: code = Canceled desc = oopsie",
			},
		},
		{
			name:    "Canceled",
			service: "foo",
			err:     status.Errorf(codes.Unimplemented, "oopsie"),

			expectLogLevel: loggers.LogLevelWarning,
			expectConsole:  regexp.MustCompile("^⟁ Unimplemented \\[foo]\n\n  - rpc error: code = Unimplemented desc = oopsie\n\n$"),
			expectJSON: map[string]interface{}{
				"grpcRequest": map[string]interface{}{
					"code":    codes.Unimplemented,
					"service": "foo",
				},
				"severity": "WARNING",
				"error":    "rpc error: code = Unimplemented desc = oopsie",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formatter := formattersmocks.NewMockFormatter(t)
			adapter := adapters.NewGRPC(formatter)

			var consoleResult string
			var jsonResult interface{}

			formatter.
				On("Log", mock.Anything, tc.expectLogLevel).
				Run(func(args mock.Arguments) {
					content := args.Get(0).(formatters.LogSplit)
					consoleResult = content.RenderConsole()
					jsonResult = content.RenderJSON()
				})

			adapter.Report(tc.service, tc.err)

			require.Regexp(t, tc.expectConsole, consoleResult)
			require.Equal(t, tc.expectJSON, jsonResult)
		})
	}
}

func TestWrapGRPCCall(t *testing.T) {
	t.Run("WithDefaultImplementation", func(t *testing.T) {
		formatter := formattersmocks.NewMockFormatter(t)
		adapter := adapters.NewGRPC(formatter)

		var consoleResult string
		var jsonResult interface{}

		formatter.
			On("Log", mock.Anything, loggers.LogLevelInfo).
			Run(func(args mock.Arguments) {
				content := args.Get(0).(formatters.LogSplit)
				consoleResult = content.RenderConsole()
				jsonResult = content.RenderJSON()
			})

		callback := func(ctx context.Context, in int) (int, error) {
			return in, nil
		}

		wrapped := adapters.WrapGRPCCall("foo", adapter, callback)
		out, err := wrapped(context.Background(), 123456)

		require.NoError(t, err)
		require.Equal(t, 123456, out)

		expectConsole := regexp.MustCompile("^✓ OK \\[foo] \\([^)]+\\)\n\n$")
		expectJSON := map[string]interface{}{
			"grpcRequest": map[string]interface{}{
				"code":    codes.OK,
				"service": "foo",
			},
			"severity": "INFO",
		}

		require.Regexp(t, expectConsole, consoleResult)
		require.Empty(t, cmp.Diff(
			expectJSON, jsonResult,
			cmpopts.IgnoreMapEntries(func(k string, _ interface{}) bool {
				return k == "latency"
			}),
		))

		formatter.AssertExpectations(t)
	})

	t.Run("WithAnyImplementation", func(t *testing.T) {
		adapter := adaptersmocks.NewMockGRPC(t)

		adapter.On("Report", "foo", nil).Once()

		callback := func(ctx context.Context, in int) (int, error) {
			return in, nil
		}

		wrapped := adapters.WrapGRPCCall("foo", adapter, callback)
		out, err := wrapped(context.Background(), 123456)

		require.NoError(t, err)
		require.Equal(t, 123456, out)

		adapter.AssertExpectations(t)
	})
}

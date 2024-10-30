package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/a-novel/golib/loggers/formatters"
)

func TestDiscoverGRPC(t *testing.T) {
	restore := configureTerminal()
	defer restore()

	t.Run("Render", func(t *testing.T) {
		content := formatters.NewDiscoverGRPC([]grpc.ServiceDesc{
			{
				ServiceName: "Service1",
				Methods: []grpc.MethodDesc{
					{
						MethodName: "Method1",
					},
					{
						MethodName: "Method2",
					},
				},
				Streams: []grpc.StreamDesc{
					{
						StreamName: "Stream1",
					},
					{
						StreamName: "Stream2",
					},
				},
			},
			{
				ServiceName: "Service2",
				Methods: []grpc.MethodDesc{
					{
						MethodName: "Method1",
					},
				},
			},
			{
				ServiceName: "Service3",
				Methods: []grpc.MethodDesc{
					{
						MethodName: "Method1",
					},
					{
						MethodName: "Method2",
					},
				},
			},
		}, 1234)

		expectConsole := "\x1b[38;2;0;167;255m╭────────────────────────────────────────────────────────────────╮\x1b[0m\n" +
			"\x1b[38;2;0;167;255m│\x1b[0m \x1b[1;38;2;0;167;255mRPC services successfully registered.\x1b[0m                          \x1b[38;2;0;167;255m│\x1b[0m\n" +
			"\x1b[38;2;0;167;255m│\x1b[0m \x1b[2;38;2;0;167;255m3 services registered on port :1234\x1b[0m                            \x1b[38;2;0;167;255m│\x1b[0m\n" +
			"\x1b[38;2;0;167;255m╰────────────────────────────────────────────────────────────────╯\x1b[0m\n\n" +
			"     \x1b[38;2;255;204;102mService1\x1b[0m\n" +
			"     \x1b[38;2;255;204;102m    \x1b[2;38;2;255;204;102mMethod1\x1b[0m\x1b[0m  \n" +
			"     \x1b[38;2;255;204;102m    \x1b[2;38;2;255;204;102mMethod2\x1b[0m\x1b[0m  \n" +
			"     \x1b[38;2;255;204;102m    \x1b[2;38;2;255;204;102m[Stream1]\x1b[0m\x1b[0m\n" +
			"     \x1b[38;2;255;204;102m    \x1b[2;38;2;255;204;102m[Stream2]\x1b[0m\x1b[0m\n" +
			"     \x1b[38;2;255;204;102mService2\x1b[0m\n" +
			"     \x1b[38;2;255;204;102m    \x1b[2;38;2;255;204;102mMethod1\x1b[0m\x1b[0m\n" +
			"     \x1b[38;2;255;204;102mService3\x1b[0m\n" +
			"     \x1b[38;2;255;204;102m    \x1b[2;38;2;255;204;102mMethod1\x1b[0m\x1b[0m\n" +
			"     \x1b[38;2;255;204;102m    \x1b[2;38;2;255;204;102mMethod2\x1b[0m\x1b[0m\n\n"
		expectJSON := map[string]interface{}{
			"Service1": map[string]interface{}{
				"methods": []interface{}{"Method1", "Method2"},
				"streams": []interface{}{"Stream1", "Stream2"},
			},
			"Service2": map[string]interface{}{
				"methods": []interface{}{"Method1"},
				"streams": []interface{}{},
			},
			"Service3": map[string]interface{}{
				"methods": []interface{}{"Method1", "Method2"},
				"streams": []interface{}{},
			},
		}

		require.Equal(t, expectConsole, content.RenderConsole())
		require.Equal(t, expectJSON, content.RenderJSON())
	})
}

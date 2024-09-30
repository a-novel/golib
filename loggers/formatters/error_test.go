package formatters_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/loggers/formatters"
)

func TestError(t *testing.T) {
	restore := configureTerminal()
	defer restore()

	content := formatters.NewError(errors.New("foo bar qux"), "nawak")

	require.Equal(t, "\x1b[38;2;255;50;50mnawak: foo bar qux\x1b[0m\n", content.RenderConsole())
	require.Equal(t, map[string]interface{}{"error": "foo bar qux", "message": "nawak"}, content.RenderJSON())
}

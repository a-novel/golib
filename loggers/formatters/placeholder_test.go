package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/loggers/formatters"
)

func TestPlaceholder(t *testing.T) {
	restore := configureTerminal()
	defer restore()

	content := formatters.NewPlaceholder("foo bar qux")

	require.Equal(t, "\x1b[38;2;255;128;0m⚠  \x1b[0m\x1b[2mfoo bar qux\x1b[0m\n", content.RenderConsole())
	require.Equal(t, map[string]interface{}{"message": "foo bar qux"}, content.RenderJSON())
}

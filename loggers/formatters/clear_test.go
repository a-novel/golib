package formatters_test

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/loggers/formatters"
)

func TestClear(t *testing.T) {
	restore := configureTerminal()
	defer restore()

	content := formatters.NewClear()

	require.Equal(t, ansi.EraseDisplay(2)+ansi.CursorOrigin, content.RenderConsole())
	require.Nil(t, content.RenderJSON())
}

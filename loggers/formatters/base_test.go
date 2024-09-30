package formatters_test

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/loggers/formatters"
)

func TestLogBase(t *testing.T) {
	t.Run("Render", func(t *testing.T) {
		restore := configureTerminal()
		defer restore()

		content := formatters.NewBase("foo bar qux")

		require.Equal(t, "foo bar qux\n", content.RenderConsole())
		require.Equal(t, map[string]interface{}{"message": "foo bar qux"}, content.RenderJSON())
	})

	t.Run("Style", func(t *testing.T) {
		restore := configureTerminal()
		defer restore()

		content := formatters.NewBase("foo bar qux").
			SetStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8000")))

		require.Equal(t, "\x1b[38;2;255;128;0mfoo bar qux\x1b[0m\n", content.RenderConsole())
		require.Equal(t, map[string]interface{}{"message": "foo bar qux"}, content.RenderJSON())
	})

	t.Run("Child", func(t *testing.T) {
		restore := configureTerminal()
		defer restore()

		content := formatters.NewBase("foo bar qux").
			SetStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8000")))
		child := formatters.NewBase("child message")

		content.SetChild(child)

		require.Equal(t, "\x1b[38;2;255;128;0mfoo bar qux\x1b[0m\nchild message\n", content.RenderConsole())
		require.Equal(
			t,
			map[string]interface{}{
				"message": "foo bar qux",
				"child":   map[string]interface{}{"message": "child message"},
			},
			content.RenderJSON(),
		)
	})
}

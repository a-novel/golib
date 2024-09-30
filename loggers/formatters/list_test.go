package formatters_test

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/loggers/formatters"
)

func TestList(t *testing.T) {
	restore := configureTerminal()
	defer restore()

	content := formatters.NewList().
		Append("foo").
		Append("bar").
		Append("qux")

	require.Equal(t, "foo\nbar\nqux\n", content.RenderConsole())
	require.Equal(t, []interface{}{"foo", "bar", "qux"}, content.RenderJSON())

	t.Run("NoItems", func(t *testing.T) {
		content := formatters.NewList()
		require.Equal(t, "", content.RenderConsole())
		require.Equal(t, nil, content.RenderJSON())
	})

	t.Run("Nested", func(t *testing.T) {
		content := formatters.NewList().
			Append("foo").
			Nest(
				formatters.NewList().
					Append("bar").
					Nest(
						formatters.NewBase("baz").SetStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3232"))),
					),
			).
			Append("qux")

		require.Equal(t, "foo\n bar\n \x1b[38;2;255;50;50mbaz\x1b[0m\nqux\n", content.RenderConsole())
		require.Equal(t, []interface{}{
			"foo",
			[]interface{}{
				"bar",
				map[string]interface{}{"message": "baz"},
			},
			"qux",
		}, content.RenderJSON())
	})

	t.Run("Customized", func(t *testing.T) {
		content := formatters.NewList().
			Append("foo").
			Nest(
				formatters.NewList().
					Append("bar"),
			).
			Append("qux")

		content.
			SetEnumerator(list.Dash).
			SetIndenter(func(items list.Items, index int) string {
				return "-->"
			}).
			SetStyle(lipgloss.NewStyle().Bold(true)).
			SetItemStyle(lipgloss.NewStyle().Faint(true))

		require.Equal(
			t,
			"\x1b[1m-\x1b[0m\x1b[2mfoo\x1b[0m\n\x1b[1m-\x1b[0m\x1b[2mbar\x1b[0m\n\x1b[1m-\x1b[0m\x1b[2mqux\x1b[0m\n",
			content.RenderConsole(),
		)
		require.Equal(t, []interface{}{
			"foo",
			[]interface{}{"bar"},
			"qux",
		}, content.RenderJSON())
	})
}

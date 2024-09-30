package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/loggers/formatters"
)

func TestSplit(t *testing.T) {
	t.Run("Render", func(t *testing.T) {
		content := formatters.NewSplit()

		content.SetConsoleContent(formatters.NewBase("foo"))
		content.SetJSONContent(formatters.NewBase("bar"))

		require.Equal(t, "foo\n", content.RenderConsole())
		require.Equal(t, map[string]interface{}{"message": "bar"}, content.RenderJSON())
	})

	t.Run("RenderMessages", func(t *testing.T) {
		content := formatters.NewSplit()

		content.SetConsoleMessage("foo")
		content.SetJSONMessage("bar")

		require.Equal(t, "foo\n", content.RenderConsole())
		require.Equal(t, "bar", content.RenderJSON())
	})

	t.Run("RenderFns", func(t *testing.T) {
		content := formatters.NewSplit()

		content.SetConsoleRenderer(func() string { return "foo" })
		content.SetJSONRenderer(func() interface{} { return "bar" })

		require.Equal(t, "foo", content.RenderConsole())
		require.Equal(t, "bar", content.RenderJSON())
	})

	t.Run("ConsoleOnly", func(t *testing.T) {
		content := formatters.NewSplit()

		content.SetConsoleContent(formatters.NewBase("foo"))

		require.Equal(t, "foo\n", content.RenderConsole())
		require.Equal(t, nil, content.RenderJSON())
	})

	t.Run("JSONOnly", func(t *testing.T) {
		content := formatters.NewSplit()

		content.SetJSONContent(formatters.NewBase("bar"))

		require.Equal(t, "", content.RenderConsole())
		require.Equal(t, map[string]interface{}{"message": "bar"}, content.RenderJSON())
	})

	t.Run("Nothing", func(t *testing.T) {
		content := formatters.NewSplit()

		require.Equal(t, "", content.RenderConsole())
		require.Equal(t, nil, content.RenderJSON())
	})
}

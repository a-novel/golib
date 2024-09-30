package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/loggers/formatters"
)

func TestTitle(t *testing.T) {
	t.Run("Render", func(t *testing.T) {
		restore := configureTerminal()
		defer restore()

		content := formatters.NewTitle("UmU")

		expectConsole := "\x1b[38;2;0;167;255m╭────────────────────────────────────────────────────────────────╮\x1b[0m\n\x1b[38;2;0;167;255m│\x1b[0m \x1b[1;38;2;0;167;255mUmU\x1b[0m                                                            \x1b[38;2;0;167;255m│\x1b[0m\n\x1b[38;2;0;167;255m╰────────────────────────────────────────────────────────────────╯\x1b[0m\n\n"

		require.Equal(t, expectConsole, content.RenderConsole())
		require.Equal(t, map[string]interface{}{"message": "UmU"}, content.RenderJSON())
	})

	t.Run("Description", func(t *testing.T) {
		restore := configureTerminal()
		defer restore()

		content := formatters.NewTitle("UmU").SetDescription("uwu omo owo")

		expectConsole := "\x1b[38;2;0;167;255m╭────────────────────────────────────────────────────────────────╮\x1b[0m\n\x1b[38;2;0;167;255m│\x1b[0m \x1b[1;38;2;0;167;255mUmU\x1b[0m                                                            \x1b[38;2;0;167;255m│\x1b[0m\n\x1b[38;2;0;167;255m│\x1b[0m \x1b[2;38;2;0;167;255muwu omo owo\x1b[0m                                                    \x1b[38;2;0;167;255m│\x1b[0m\n\x1b[38;2;0;167;255m╰────────────────────────────────────────────────────────────────╯\x1b[0m\n\n"

		require.Equal(t, expectConsole, content.RenderConsole())
		require.Equal(t, map[string]interface{}{"message": "UmU", "description": "uwu omo owo"}, content.RenderJSON())
	})

	t.Run("Child", func(t *testing.T) {
		restore := configureTerminal()
		defer restore()

		content := formatters.NewTitle("UmU").SetChild(formatters.NewBase("owumo"))

		expectConsole := "\x1b[38;2;0;167;255m╭────────────────────────────────────────────────────────────────╮\x1b[0m\n\x1b[38;2;0;167;255m│\x1b[0m \x1b[1;38;2;0;167;255mUmU\x1b[0m                                                            \x1b[38;2;0;167;255m│\x1b[0m\n\x1b[38;2;0;167;255m╰────────────────────────────────────────────────────────────────╯\x1b[0m\n\nowumo\n\n"

		require.Equal(t, expectConsole, content.RenderConsole())
		require.Equal(
			t,
			map[string]interface{}{"message": "UmU", "data": map[string]interface{}{"message": "owumo"}},
			content.RenderJSON(),
		)
	})
}

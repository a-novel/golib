package formatters_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/loggers/formatters"
)

var dummySpinner = spinner.Spinner{
	Frames: []string{"O", "w", "O"},
	FPS:    100 * time.Millisecond,
}

func requireLoaderConsole(t *testing.T, expect string, content formatters.LogLoader) {
	// There is a time elapsed string at the end that we ignore, since it is impossible to predict.
	rendered := content.RenderConsole()
	require.True(
		t,
		strings.HasPrefix(rendered, expect),
		fmt.Sprintf("got: %s\nexpected prefix: %s", escapeAnsi(rendered), escapeAnsi(expect)),
	)
}

func requireLoaderConsoleChild(t *testing.T, expect string, content formatters.LogLoader) {
	// There is a time elapsed string at the end that we ignore, since it is impossible to predict.
	rendered := content.RenderConsole()
	require.True(
		t,
		strings.HasSuffix(rendered, expect),
		fmt.Sprintf("got: %s\nexpected suffix: %s", escapeAnsi(rendered), escapeAnsi(expect)),
	)
}

func requireLoaderJSON(t *testing.T, expect map[string]interface{}, content formatters.LogLoader) {
	rendered, ok := content.RenderJSON().(map[string]interface{})
	require.True(t, ok)

	elapsedRaw, ok := rendered["elapsed"]
	require.True(t, ok)
	require.NotEmpty(t, elapsedRaw)
	delete(rendered, "elapsed")

	elapsedNanos, ok := rendered["elapsed_nanos"]
	require.True(t, ok)
	require.NotEmpty(t, elapsedNanos)
	delete(rendered, "elapsed_nanos")

	opID, ok := rendered["op_id"]
	require.True(t, ok)
	require.NotEmpty(t, opID)
	delete(rendered, "op_id")

	require.Equal(t, expect, rendered)
}

func TestLoader(t *testing.T) {
	t.Run("Static", func(t *testing.T) {
		restore := configureTerminal()
		defer restore()

		content := formatters.NewLoader("foo bar qux", dummySpinner)

		requireLoaderConsole(t, "\x1b[38;2;255;0;127mO\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)

		time.Sleep(110 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mO\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)
	})

	t.Run("Dynamic", func(t *testing.T) {
		restore := configureTerminal()
		defer restore()

		content := formatters.NewLoader("foo bar qux", dummySpinner)
		renderer := newDynamicRenderer()
		content.RenderConsoleDynamic(renderer.renderer)

		time.Sleep(10 * time.Millisecond)

		requireLoaderConsole(t, "\x1b[38;2;255;0;127mO\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)

		time.Sleep(100 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mw\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)

		time.Sleep(100 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mO\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)

		// New cycle.
		time.Sleep(100 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mO\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)

		// Stop the loader.
		content.StopRunning()

		time.Sleep(100 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mO\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)
	})

	t.Run("UpdateDescription", func(t *testing.T) {
		restore := configureTerminal()
		defer restore()

		content := formatters.NewLoader("foo bar qux", dummySpinner)
		renderer := newDynamicRenderer()
		content.RenderConsoleDynamic(renderer.renderer)

		time.Sleep(10 * time.Millisecond)

		requireLoaderConsole(t, "\x1b[38;2;255;0;127mO\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)

		content.SetDescription("nawak")
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mO\x1b[0m\x1b[2m nawak", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "nawak"}, content)

		time.Sleep(100 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mw\x1b[0m\x1b[2m nawak", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "nawak"}, content)

		content.StopRunning()

		time.Sleep(100 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mw\x1b[0m\x1b[2m nawak", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "nawak"}, content)

		content.SetDescription("foo bar qux")
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mw\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)
	})

	t.Run("SetChild", func(t *testing.T) {
		restore := configureTerminal()
		defer restore()

		content := formatters.NewLoader("foo bar qux", dummySpinner)
		renderer := newDynamicRenderer()
		content.RenderConsoleDynamic(renderer.renderer)

		time.Sleep(10 * time.Millisecond)

		requireLoaderConsole(t, "\x1b[38;2;255;0;127mO\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)

		content.SetChild(formatters.NewBase("nawak").SetStyle(lipgloss.NewStyle().Bold(true)))
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mO\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderConsoleChild(t, "\n\x1b[1mnawak\x1b[0m\n", content)
		requireLoaderJSON(
			t,
			map[string]interface{}{
				"message": "foo bar qux",
				"data":    map[string]interface{}{"message": "nawak"},
			},
			content,
		)

		time.Sleep(100 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mw\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderConsoleChild(t, "\n\x1b[1mnawak\x1b[0m\n", content)
		requireLoaderJSON(
			t,
			map[string]interface{}{
				"message": "foo bar qux",
				"data":    map[string]interface{}{"message": "nawak"},
			},
			content,
		)

		content.StopRunning()

		time.Sleep(100 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mw\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderConsoleChild(t, "\n\x1b[1mnawak\x1b[0m\n", content)
		requireLoaderJSON(
			t,
			map[string]interface{}{
				"message": "foo bar qux",
				"data":    map[string]interface{}{"message": "nawak"},
			},
			content,
		)

		content.SetChild(nil)
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mw\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)
	})

	t.Run("Completed", func(t *testing.T) {
		restore := configureTerminal()
		defer restore()

		content := formatters.NewLoader("foo bar qux", dummySpinner)
		renderer := newDynamicRenderer()
		content.RenderConsoleDynamic(renderer.renderer)

		time.Sleep(110 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mw\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)

		content.SetCompleted()
		// The timer constantly creates new messages, so comparing them one by one is impossible.
		// Instead, we check no new messages are sent, to ensure the loader is stopped.
		messagesCount := len(renderer.getCalls())
		requireLoaderConsole(t, "\x1b[38;2;0;255;0m✓\x1b[0m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux", "completed": true}, content)

		time.Sleep(110 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;0;255;0m✓\x1b[0m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux", "completed": true}, content)
		require.Equal(t, messagesCount, len(renderer.getCalls()))
	})

	t.Run("Error", func(t *testing.T) {
		restore := configureTerminal()
		defer restore()

		content := formatters.NewLoader("foo bar qux", dummySpinner)
		renderer := newDynamicRenderer()
		content.RenderConsoleDynamic(renderer.renderer)

		time.Sleep(110 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;255;0;127mw\x1b[0m\x1b[2m foo bar qux", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux"}, content)

		content.SetError()
		// The timer constantly creates new messages, so comparing them one by one is impossible.
		// Instead, we check no new messages are sent, to ensure the loader is stopped.
		messagesCount := len(renderer.getCalls())
		requireLoaderConsole(t, "\x1b[38;2;255;50;50m✗ foo bar qux\x1b[0m\x1b[2m", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux", "error": true}, content)

		time.Sleep(110 * time.Millisecond)
		requireLoaderConsole(t, "\x1b[38;2;255;50;50m✗ foo bar qux\x1b[0m\x1b[2m", content)
		requireLoaderJSON(t, map[string]interface{}{"message": "foo bar qux", "error": true}, content)
		require.Equal(t, messagesCount, len(renderer.getCalls()))
	})
}

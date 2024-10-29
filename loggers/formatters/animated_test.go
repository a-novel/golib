package formatters_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/loggers/formatters"
)

// Emulates a rendering function for dynamic content. Each time a render is triggered, the rendered content is
// captured inside the calls slice.
type dummyDynamicRenderer struct {
	// Each time a new render is performed, the message is appended to this slice.
	calls []string
}

// The main render method, just capture the rendered message into the calls slice.
func (d *dummyDynamicRenderer) renderer(msg string) {
	d.calls = append(d.calls, msg)
}

// Returns the results of each call to the renderer.
func (d *dummyDynamicRenderer) getCalls() []string {
	return d.calls
}

// The most basic dynamic rendered. It just keeps track of the number of times it was updated.
type dummyDynamicLogContent struct {
	// Circular reference that implements all the dynamic logic.
	formatters.LogDynamicContentPure
	// Updates tracker.
	updatedTimes int
}

// Increment the update tracker.
func (l *dummyDynamicLogContent) update() {
	l.updatedTimes++
}

func (l *dummyDynamicLogContent) RenderConsole() string {
	return fmt.Sprintf("Updated %d times", l.updatedTimes)
}

func (l *dummyDynamicLogContent) RenderJSON() interface{} {
	return map[string]interface{}{"updated_times": l.updatedTimes}
}

// Returns a new instance of the dummy dynamic log content, for testing purposes.
func newDummyDynamicLogContent() *dummyDynamicLogContent {
	content := &dummyDynamicLogContent{}
	content.LogDynamicContentPure = formatters.NewAnimated(content, 100*time.Millisecond, content.update)
	return content
}

// Returns a new instance of the dummy dynamic renderer, for testing purposes.
func newDynamicRenderer() *dummyDynamicRenderer {
	return &dummyDynamicRenderer{calls: []string{}}
}

// Assert the rendered value of the dummy dynamic log content.
func requireDynamicLogContent(t *testing.T, content *dummyDynamicLogContent, updates int) {
	require.Equal(t, fmt.Sprintf("Updated %v times", updates), content.RenderConsole())
	require.Equal(t, map[string]interface{}{"updated_times": updates}, content.RenderJSON())
}

// Assert the rendered value of the dummy dynamic log content, with a dynamic renderer.
func requireDynamicLogContentRenders(t *testing.T, renderer *dummyDynamicRenderer, updates int) {
	require.Equal(t, updates+1, len(renderer.getCalls()))
	for i := 0; i < updates; i++ {
		require.Equal(t, fmt.Sprintf(formatters.EraseLineSequence+"Updated %v times", i), renderer.getCalls()[i], i)
	}
}

func TestLogAnimated(t *testing.T) {
	t.Run("Static", func(t *testing.T) {
		t.Parallel()

		content := newDummyDynamicLogContent()

		requireDynamicLogContent(t, content, 0)

		// Content should remain the same no matter how long we wait.
		time.Sleep(120 * time.Millisecond)
		requireDynamicLogContent(t, content, 0)
	})

	t.Run("Dynamic", func(t *testing.T) {
		t.Parallel()

		content := newDummyDynamicLogContent()

		renderer := newDynamicRenderer()
		content.RenderConsoleDynamic(renderer.renderer)

		time.Sleep(10 * time.Millisecond)

		requireDynamicLogContentRenders(t, renderer, 0)
		requireDynamicLogContent(t, content, 0)

		// Test content update.
		time.Sleep(100 * time.Millisecond)

		requireDynamicLogContentRenders(t, renderer, 1)
		requireDynamicLogContent(t, content, 1)

		time.Sleep(100 * time.Millisecond)

		requireDynamicLogContentRenders(t, renderer, 2)
		requireDynamicLogContent(t, content, 2)

		// Test content stop.
		content.StopRunning()

		requireDynamicLogContentRenders(t, renderer, 2)
		requireDynamicLogContent(t, content, 2)

		time.Sleep(200 * time.Millisecond)

		requireDynamicLogContentRenders(t, renderer, 2)
		requireDynamicLogContent(t, content, 2)
	})

	t.Run("StopUnstarted", func(t *testing.T) {
		t.Parallel()

		content := newDummyDynamicLogContent()

		renderer := newDynamicRenderer()
		content.StopRunning()

		requireDynamicLogContent(t, content, 0)

		content.RenderConsoleDynamic(renderer.renderer)

		time.Sleep(210 * time.Millisecond)

		requireDynamicLogContentRenders(t, renderer, 2)
		requireDynamicLogContent(t, content, 2)

		// Test content stop.
		content.StopRunning()

		time.Sleep(200 * time.Millisecond)

		requireDynamicLogContentRenders(t, renderer, 2)
		requireDynamicLogContent(t, content, 2)
	})

	t.Run("NoRestart", func(t *testing.T) {
		t.Parallel()

		content := newDummyDynamicLogContent()

		renderer := newDynamicRenderer()
		content.RenderConsoleDynamic(renderer.renderer)

		// Test content update.
		time.Sleep(110 * time.Millisecond)

		requireDynamicLogContentRenders(t, renderer, 1)
		requireDynamicLogContent(t, content, 1)

		content.StopRunning()
		time.Sleep(100 * time.Millisecond)

		requireDynamicLogContentRenders(t, renderer, 1)
		requireDynamicLogContent(t, content, 1)

		content.RenderConsoleDynamic(renderer.renderer)
		time.Sleep(210 * time.Millisecond)

		requireDynamicLogContentRenders(t, renderer, 1)
		requireDynamicLogContent(t, content, 1)
	})

	t.Run("Concurrency/RenderDynamic", func(t *testing.T) {
		t.Parallel()

		content := newDummyDynamicLogContent()

		renderer := newDynamicRenderer()
		// This renderer will capture messages from concurrent calls to RenderConsole while the dynamic content is
		// running.
		concurrentRenderer := newDynamicRenderer()
		content.RenderConsoleDynamic(renderer.renderer)

		// Test content update.
		time.Sleep(110 * time.Millisecond)

		requireDynamicLogContentRenders(t, renderer, 1)
		requireDynamicLogContent(t, content, 1)

		wg := new(sync.WaitGroup)

		callInitialization := func(i int) {
			defer wg.Done()
			content.RenderConsoleDynamic(renderer.renderer)
			concurrentRenderer.renderer(content.RenderConsole())
		}

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go callInitialization(i)
		}

		wg.Wait()
		content.StopRunning()

		requireDynamicLogContentRenders(t, renderer, 1)
		requireDynamicLogContent(t, content, 1)

		require.Len(t, concurrentRenderer.getCalls(), 10)
		for i := 0; i < 10; i++ {
			require.Equal(t, "Updated 1 times", concurrentRenderer.getCalls()[i])
		}
	})

	t.Run("Concurrency/StopRunning", func(t *testing.T) {
		t.Parallel()

		content := newDummyDynamicLogContent()

		renderer := newDynamicRenderer()
		// This renderer will capture messages from concurrent calls to RenderConsole while the dynamic content is
		// running.
		concurrentRenderer := newDynamicRenderer()
		content.RenderConsoleDynamic(renderer.renderer)

		// Test content update.
		time.Sleep(110 * time.Millisecond)

		requireDynamicLogContentRenders(t, renderer, 1)
		requireDynamicLogContent(t, content, 1)

		wg := new(sync.WaitGroup)

		callStopRunning := func(i int) {
			defer wg.Done()
			content.StopRunning()
			concurrentRenderer.renderer(content.RenderConsole())
		}

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go callStopRunning(i)
		}

		wg.Wait()

		time.Sleep(100 * time.Millisecond)

		requireDynamicLogContentRenders(t, renderer, 1)
		requireDynamicLogContent(t, content, 1)
	})
}

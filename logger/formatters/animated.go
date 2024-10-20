package formatters

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"sync"
	"time"
)

// LogDynamicContentPure is a subset of LogDynamicContent implemented by the generic default implementation.
type LogDynamicContentPure interface {
	// RenderConsoleDynamic implements LogDynamicContent.RenderConsoleDynamic interface.
	RenderConsoleDynamic(renderer func(msg string)) LogDynamicContent
	// StopRunning implements LogDynamicContent.StopRunning interface.
	StopRunning() LogDynamicContent
}

// Default implementation of LogDynamicContentPure, allowing to make custom LogContent dynamic, without rewriting
// the whole dynamic logic.
//
// It allows the message to contain animated content, that automatically updates independently of the message itself.
type animatedLogContentImpl struct {
	// A circular reference to the source content this implementation extends.
	parentContent LogDynamicContent

	// The function, passed by the Formatter, that should be used to print new content.
	renderFn func(msg string)
	// FPS is the minimum rate at which the animated content should be updated.
	fps time.Duration
	// A function to update the animated content. It is the parentContent responsibility to refresh the render of
	// the animated content when this is triggered.
	update func()

	// Holds the termination signal.
	quit chan bool

	// Keep track of the last time the update method was called, so it can be properly triggered depending on the fps.
	lastRendered time.Time
	// A slight optimization, that allows to skip updates if the rendered content does not change.
	lastRenderedValue string

	// Prevent race conditions.
	mu sync.Mutex
	// Allows the StopRunning method to wait for the rendering goroutine to finish. Otherwise, this instance can be
	// deallocated before the message has time to be rendered.
	wg sync.WaitGroup
}

// Trigger a repainting of the content.
func (l *animatedLogContentImpl) exec() {
	// Force only 1 render at a time.
	l.mu.Lock()

	// Check if the animated content can be updated.
	canRender := time.Since(l.lastRendered) > l.fps
	if canRender {
		// Trigger the animated content update.
		l.update()
		// Update the last rendered time.
		l.lastRendered = time.Now()
	}

	// Get the new rendered value from the parent.
	newRender := l.parentContent.RenderConsole()
	// Don't do anything if the output has not changed.
	if newRender == l.lastRenderedValue {
		l.mu.Unlock()
		return
	}

	// Overwrite previous content.
	lastRenderedHeight := lipgloss.Height(l.lastRenderedValue)
	for i := 0; i < lastRenderedHeight; i++ {
		l.renderFn(ansi.EraseEntireLine + "\r" + ansi.CursorUp1)
	}

	l.lastRenderedValue = newRender
	l.mu.Unlock()

	// Send the new content to the renderer.
	l.renderFn(newRender)
}

// The main goroutine used to render the animated content.
func (l *animatedLogContentImpl) renderConsoleDynamic() {
	// Properly mark the goroutine as finished.
	defer l.wg.Done()
	for {
		select {
		case <-l.quit:
			return
		default:
			l.exec()
		}
	}
}

// Returns whether the dynamic loop is running.
func (l *animatedLogContentImpl) isRunning() bool {
	if l.quit == nil {
		return false
	}

	quitted, closed := <-l.quit
	return !quitted && !closed
}

// Returns whether the dynamic loop is idle (not started yet).
func (l *animatedLogContentImpl) isIdle() bool {
	return l.quit == nil
}

// Returns whether the dynamic loop has terminated.
func (l *animatedLogContentImpl) isTerminated() bool {
	if l.quit == nil {
		return false
	}

	_, closed := <-l.quit
	return closed
}

// RenderConsoleDynamic implements LogDynamicContent.RenderConsoleDynamic interface.
func (l *animatedLogContentImpl) RenderConsoleDynamic(rendererFn func(msg string)) LogDynamicContent {
	l.mu.Lock()
	if !l.isIdle() {
		// If the dynamic loop is still running, update the rendering method.
		if l.isRunning() {
			l.renderFn = rendererFn
		}

		l.mu.Unlock()
		// If the dynamic loop has terminated, this is a no-op.
		return l.parentContent
	}

	// Initialize the dynamic loop.
	l.renderFn = rendererFn
	l.quit = make(chan bool)
	l.wg.Add(1)
	l.mu.Unlock()
	// Run the dynamic loop.
	go l.renderConsoleDynamic()
	return l.parentContent
}

// StopRunning implements LogDynamicContent.StopRunning interface.
func (l *animatedLogContentImpl) StopRunning() LogDynamicContent {
	// Loop must have been initialized and not terminated yet, for this method to actually do something.
	if l.isRunning() {
		// Trigger the termination of the dynamic loop.
		l.quit <- true
		l.wg.Wait()

		l.mu.Lock()
		// Definitely close the channel. It cannot be re-opened past this point.
		close(l.quit)
		l.mu.Unlock()

		// Trigger one last render, to make sure the content is up-to-date.
		l.exec()

		// Deallocate the rendering function.
		l.mu.Lock()
		l.renderFn = nil
		l.mu.Unlock()
	}

	return l.parentContent
}

// NewAnimatedLogContent creates a new LogDynamicContentPure instance, that allows to render animated content.
func NewAnimatedLogContent(renderer LogDynamicContent, fps time.Duration, update func()) LogDynamicContentPure {
	return &animatedLogContentImpl{
		parentContent: renderer,
		fps:           fps,
		update:        update,
	}
}

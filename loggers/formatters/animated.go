package formatters

import (
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

const EraseLineSequence = ansi.EraseEntireLine + "\r" + ansi.CursorUp1

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
	// Set to true when the channel is closed, intended to prevent readings from a closed channel.
	terminated bool

	// Keep track of the last time the update method was called, so it can be properly triggered depending on the fps.
	lastRendered time.Time
	// A slight optimization, that allows to skip updates if the rendered content does not change.
	lastRenderedValue string

	// Allows the StopRunning method to wait for the rendering goroutine to finish. Otherwise, this instance can be
	// deallocated before the message has time to be rendered.
	wg sync.WaitGroup
	mu sync.Mutex
}

// Trigger a repainting of the content.
func (animatedLogContent *animatedLogContentImpl) exec() {
	animatedLogContent.mu.Lock()
	lastRRendered := animatedLogContent.lastRendered
	animatedLogContent.mu.Unlock()

	// Check if the animated content can be updated.
	shouldUpdateAnimation := time.Since(lastRRendered) > animatedLogContent.fps

	if shouldUpdateAnimation {
		// Only trigger update once the render has been updated once, otherwise keep the first frame.
		if lastRRendered.After(epoch) {
			// Trigger the animated content update.
			// We don't call mutex here, because parent update method should be thread-safe.
			animatedLogContent.update()
		}

		// Update the last rendered time.
		animatedLogContent.mu.Lock()
		animatedLogContent.lastRendered = time.Now()
		animatedLogContent.mu.Unlock()
	}

	// Get the new rendered value from the parent.
	newRender := animatedLogContent.parentContent.RenderConsole()

	// From now on, we enter the rendering process. We want to prevent concurrent renders.
	animatedLogContent.mu.Lock()

	// Don't do anything if the output has not changed.
	if newRender == animatedLogContent.lastRenderedValue {
		animatedLogContent.mu.Unlock()
		return
	}

	animatedLogContent.lastRenderedValue = newRender

	// Overwrite previous content.
	for range lipgloss.Height(animatedLogContent.lastRenderedValue) {
		newRender = EraseLineSequence + newRender
	}

	animatedLogContent.mu.Unlock()

	// Send the new content to the renderer.
	animatedLogContent.renderFn(newRender)
}

// The main goroutine used to render the animated content.
func (animatedLogContent *animatedLogContentImpl) renderConsoleDynamic() {
	// Properly mark the goroutine as finished.
	defer animatedLogContent.wg.Done()
	for {
		select {
		case <-animatedLogContent.quit:
			return
		default:
			animatedLogContent.exec()
		}
	}
}

// RenderConsoleDynamic implements LogDynamicContent.RenderConsoleDynamic interface.
func (animatedLogContent *animatedLogContentImpl) RenderConsoleDynamic(rendererFn func(msg string)) LogDynamicContent {
	// Initialize the dynamic loop.
	animatedLogContent.mu.Lock()
	animatedLogContent.renderFn = rendererFn
	animatedLogContent.mu.Unlock()

	// THe animator is idle, has not started yet.
	if animatedLogContent.quit == nil {
		animatedLogContent.mu.Lock()
		animatedLogContent.quit = make(chan bool, 1)
		animatedLogContent.wg.Add(1)
		animatedLogContent.mu.Unlock()

		// Run the dynamic loop.
		go animatedLogContent.renderConsoleDynamic()
	}

	return animatedLogContent.parentContent
}

// StopRunning implements LogDynamicContent.StopRunning interface.
func (animatedLogContent *animatedLogContentImpl) StopRunning() LogDynamicContent {
	// Loop must have been initialized and not terminated yet, for this method to actually do something.
	// Trigger the termination of the dynamic loop.
	animatedLogContent.mu.Lock()
	if animatedLogContent.terminated || animatedLogContent.quit == nil {
		animatedLogContent.mu.Unlock()
		return animatedLogContent.parentContent
	}

	animatedLogContent.terminated = true
	animatedLogContent.mu.Unlock()

	animatedLogContent.quit <- true
	animatedLogContent.wg.Wait()
	close(animatedLogContent.quit)

	// Trigger one last render, to make sure the content is up-to-date.
	// animatedLogContent.exec()

	return animatedLogContent.parentContent
}

// NewAnimated creates a new LogDynamicContentPure instance, that allows to render animated content.
func NewAnimated(renderer LogDynamicContent, fps time.Duration, update func()) LogDynamicContentPure {
	return &animatedLogContentImpl{
		parentContent: renderer,
		fps:           fps,
		update:        update,
	}
}

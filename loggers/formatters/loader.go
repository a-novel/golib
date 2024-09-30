package formatters

import (
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

// LogLoader renders a loader. Under dynamic environments, the loader is animated, as a single log. Otherwise, each
// update is rendered as a separate log.
type LogLoader interface {
	LogDynamicContent
	// SetDescription updates the main message of the loader. Each change to this value triggers a new log, under
	// static environments.
	SetDescription(description string) LogLoader
	// SetCompleted marks the loader as completed. It triggers a new log under static environments. Under dynamic
	// environments, the loader is stopped.
	SetCompleted() LogLoader
	// SetChild sets a child LogContent to the loader. Each change to this value triggers a new log, under
	// static environments.
	SetChild(child LogContent) LogLoader
	// SetError marks the loader as errored. It triggers a new log under static environments. Under dynamic
	// environments, the loader is stopped.
	SetError() LogLoader
}

// Default implementation of LogLoader.
type logLoaderImpl struct {
	LogDynamicContentPure

	isError bool

	description string
	completed   bool
	child       LogContent

	// Holds a reference to bubbles spinner.
	spinner *spinner.Model
	// Allow to display a timer, next to the description message.
	startedAt time.Time
	// Under static environments, this value links each separate log.
	opID uuid.UUID
	// Make the component thread-safe.
	mu sync.Mutex
}

// RenderConsole implements LogContent.RenderConsole interface.
func (logLoader *logLoaderImpl) RenderConsole() string {
	logLoader.mu.Lock()
	defer logLoader.mu.Unlock()

	// Compute the suffix of the message. If a child component is present, it will be rendered after the main message.
	messageSuffix := "\n"
	if logLoader.child != nil {
		messageSuffix = "\n" + logLoader.child.RenderConsole()
	}

	// Compute the time elapsed since the loader started.
	timeElapsedRaw := time.Since(logLoader.startedAt)
	// Prevent the display of values with large fractions.
	if timeElapsedRaw >= 10*time.Second {
		timeElapsedRaw = timeElapsedRaw.Round(time.Second)
	} else if timeElapsedRaw >= 10*time.Millisecond {
		timeElapsedRaw = timeElapsedRaw.Round(time.Millisecond)
	}

	// Add a space before time and description messages, so they nicely concatenate with the spinner.
	timeElapsedMessage := " " + timeElapsedRaw.String()
	descriptionMessage := " " + logLoader.description

	// The next 2 conditions render a terminated state of the loader, where the spinner is removed (since not relevant).

	if logLoader.isError {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3232")).Render("✗"+descriptionMessage) +
			lipgloss.NewStyle().Faint(true).Render(timeElapsedMessage) +
			messageSuffix
	}

	if logLoader.completed {
		ok := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("✓")
		return ok + descriptionMessage + lipgloss.NewStyle().Faint(true).Render(timeElapsedMessage) + messageSuffix
	}

	// When the loader is active, use the current frame of the spinner as a prefix to the message.
	return logLoader.spinner.View() +
		lipgloss.NewStyle().Faint(true).Render(descriptionMessage) +
		timeElapsedMessage +
		messageSuffix
}

// RenderJSON implements LogContent.RenderJSON interface.
func (logLoader *logLoaderImpl) RenderJSON() interface{} {
	output := map[string]interface{}{
		"message": logLoader.description,
		// Elapsed is preferable for human analysis, while elapsed_nanos can easily be used in queries and reporting.
		"elapsed":       time.Since(logLoader.startedAt).String(),
		"elapsed_nanos": time.Since(logLoader.startedAt).Nanoseconds(),
		// Link all messages under the same loader.
		"op_id": logLoader.opID.String(),
	}
	if logLoader.child != nil {
		output["data"] = logLoader.child.RenderJSON()
	}

	if logLoader.completed {
		output["completed"] = true
	}
	if logLoader.isError {
		output["error"] = true
	}

	return output
}

// SetDescription implements LogLoader.SetDescription interface.
func (logLoader *logLoaderImpl) SetDescription(description string) LogLoader {
	logLoader.mu.Lock()
	defer logLoader.mu.Unlock()

	logLoader.description = description
	return logLoader
}

// SetChild implements LogLoader.SetChild interface.
func (logLoader *logLoaderImpl) SetChild(child LogContent) LogLoader {
	logLoader.mu.Lock()
	defer logLoader.mu.Unlock()

	logLoader.child = child
	return logLoader
}

// SetCompleted implements LogLoader.SetCompleted interface.
func (logLoader *logLoaderImpl) SetCompleted() LogLoader {
	logLoader.mu.Lock()
	logLoader.completed = true
	logLoader.mu.Unlock()
	// Since setting this status terminates the loader, stop the spinner.
	logLoader.StopRunning()
	return logLoader
}

// SetError overrides the LogHandleError.SetError interface present in the default implementation.
func (logLoader *logLoaderImpl) SetError() LogLoader {
	logLoader.mu.Lock()
	logLoader.isError = true
	logLoader.mu.Unlock()
	// Since setting this status terminates the loader, stop the spinner.
	logLoader.StopRunning()
	return logLoader
}

// NewLoader creates a new LogLoader instance.
func NewLoader(description string, spinnerModel spinner.Spinner) LogLoader {
	// Create a new spinner from bubbles, then apply a custom style.
	spinnerInstance := spinner.New()
	spinnerInstance.Spinner = spinnerModel
	spinnerInstance.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF007F"))

	loader := &logLoaderImpl{
		spinner:     &spinnerInstance,
		description: description,
		startedAt:   time.Now(),
		opID:        uuid.New(),
	}

	// Since LogDynamicContentPure requires a circular reference to its parent, we have to set it separately.
	loader.LogDynamicContentPure = NewAnimated(loader, spinnerModel.FPS, func() {
		newSpinner, _ := spinnerInstance.Update(spinnerInstance.Tick())
		loader.mu.Lock()
		*loader.spinner = newSpinner
		loader.mu.Unlock()
	})

	return loader
}

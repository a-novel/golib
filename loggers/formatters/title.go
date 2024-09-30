package formatters

import (
	"github.com/charmbracelet/lipgloss"
)

// LogTitle renders a title element.
type LogTitle interface {
	LogContent
	// SetDescription sets the description of the title, that appears under it (inside the title block).
	SetDescription(description string) LogTitle
	// SetChild sets the child of the title, that will be rendered as a section (under the title block).
	SetChild(child LogContent) LogTitle
}

// Default implementation of LogTitle.
type logTitleImpl struct {
	title       string
	description string
	child       LogContent
}

// RenderConsole implements LogContent.RenderConsole interface.
func (logTitle *logTitleImpl) RenderConsole() string {
	titleColor := lipgloss.Color("#00A7FF")

	titleStyle := lipgloss.NewStyle().
		Foreground(titleColor).
		Bold(true).
		Align(lipgloss.Left)

	descriptionStyle := lipgloss.NewStyle().
		Foreground(titleColor).
		Bold(false).
		Faint(true).
		Align(lipgloss.Left)

	blockStyle := lipgloss.NewStyle().
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		Width(64).
		BorderForeground(titleColor)

	content := titleStyle.Render(logTitle.title)
	if logTitle.description != "" {
		content += "\n" + descriptionStyle.Render(logTitle.description)
	}

	content = blockStyle.Render(content) + "\n\n"

	if logTitle.child != nil {
		content += logTitle.child.RenderConsole() + "\n"
	}

	return content
}

// RenderJSON implements LogContent.RenderJSON interface.
func (logTitle *logTitleImpl) RenderJSON() interface{} {
	output := map[string]interface{}{"message": logTitle.title}

	if logTitle.child != nil {
		output["data"] = logTitle.child.RenderJSON()
	}

	if logTitle.description != "" {
		output["description"] = logTitle.description
	}

	return output
}

// SetDescription implements LogTitle.SetDescription interface.
func (logTitle *logTitleImpl) SetDescription(description string) LogTitle {
	logTitle.description = description
	return logTitle
}

// SetChild implements LogTitle.SetChild interface.
func (logTitle *logTitleImpl) SetChild(child LogContent) LogTitle {
	logTitle.child = child
	return logTitle
}

// NewTitle creates a new LogTitle instance.
func NewTitle(title string) LogTitle {
	return &logTitleImpl{title: title}
}

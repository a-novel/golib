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
func (l *logTitleImpl) RenderConsole() string {
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

	content := titleStyle.Render(l.title)
	if l.description != "" {
		content += "\n" + descriptionStyle.Render(l.description)
	}

	content = blockStyle.Render(content) + "\n\n"

	if l.child != nil {
		content += l.child.RenderConsole() + "\n"
	}

	return content
}

// RenderJSON implements LogContent.RenderJSON interface.
func (l *logTitleImpl) RenderJSON() interface{} {
	output := map[string]interface{}{"message": l.title}

	if l.child != nil {
		output["content"] = l.child.RenderJSON()
	}

	if l.description != "" {
		output["description"] = l.description
	}

	return output
}

// SetDescription implements LogTitle.SetDescription interface.
func (l *logTitleImpl) SetDescription(description string) LogTitle {
	l.description = description
	return l
}

// SetChild implements LogTitle.SetChild interface.
func (l *logTitleImpl) SetChild(child LogContent) LogTitle {
	l.child = child
	return l
}

// NewTitle creates a new LogTitle instance.
func NewTitle(title string) LogTitle {
	return &logTitleImpl{title: title}
}

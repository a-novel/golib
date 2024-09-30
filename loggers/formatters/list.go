package formatters

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/samber/lo"
)

// NoEnumerator adds an option to not show any enumerator on a lipgloss list.
func NoEnumerator(_ list.Items, _ int) string {
	return ""
}

// Default list indenter, no need to export as this is a lipgloss default as well.
func defaultIndenter(list.Items, int) string { return " " }

// LogList is a utility LogContent that renders a list of LogContent instances.
type LogList interface {
	LogContent
	// Append adds a new string to the list.
	Append(content string) LogList
	// Nest adds a new LogContent to the list.
	Nest(content LogContent) LogList

	// SetEnumerator sets the list enumerator.
	SetEnumerator(enumerator list.Enumerator) LogList
	// SetIndenter sets the list indenter.
	SetIndenter(indenter list.Indenter) LogList
	// SetStyle sets the style of the list.
	SetStyle(style lipgloss.Style) LogList
	// SetItemStyle sets the style of the list items.
	SetItemStyle(style lipgloss.Style) LogList
}

// Default implementation of LogList.
type logListImpl struct {
	// The ordered list of elements.
	elements []interface{}

	// ligloss list enumerator.
	enumerator list.Enumerator
	// ligloss list indenter.
	indenter list.Indenter
	// ligloss list style.
	style lipgloss.Style
	// ligloss list item style.
	itemStyle lipgloss.Style
}

// RenderConsole implements LogContent.RenderConsole interface.
func (logList *logListImpl) RenderConsole() string {
	if len(logList.elements) == 0 {
		return ""
	}

	rendered := list.New(
		lo.Map(logList.elements, func(item interface{}, _ int) interface{} {
			if renderer, ok := item.(LogContent); ok {
				// Prevent unwanted newlines.
				return strings.TrimSuffix(renderer.RenderConsole(), "\n")
			}
			return item
		})...,
	).
		Enumerator(lo.Ternary(logList.enumerator == nil, NoEnumerator, logList.enumerator)).
		Indenter(lo.Ternary(logList.indenter == nil, defaultIndenter, logList.indenter)).
		EnumeratorStyle(logList.style).
		ItemStyle(logList.itemStyle)

	return rendered.String() + "\n"
}

// RenderJSON implements LogContent.RenderJSON interface..
func (logList *logListImpl) RenderJSON() interface{} {
	if len(logList.elements) == 0 {
		return nil
	}

	return lo.Map(logList.elements, func(item interface{}, _ int) interface{} {
		if renderer, ok := item.(LogContent); ok {
			return renderer.RenderJSON()
		}
		return item
	})
}

// Append implements LogList.Append interface.
func (logList *logListImpl) Append(content string) LogList {
	logList.elements = append(logList.elements, content)
	return logList
}

// Nest implements LogList.Nest interface.
func (logList *logListImpl) Nest(content LogContent) LogList {
	logList.elements = append(logList.elements, content)
	return logList
}

// SetEnumerator implements LogList.SetEnumerator interface.
func (logList *logListImpl) SetEnumerator(enumerator list.Enumerator) LogList {
	logList.enumerator = enumerator
	return logList
}

// SetIndenter implements LogList.SetIndenter interface.
func (logList *logListImpl) SetIndenter(indenter list.Indenter) LogList {
	logList.indenter = indenter
	return logList
}

// SetStyle implements LogList.SetStyle interface.
func (logList *logListImpl) SetStyle(style lipgloss.Style) LogList {
	logList.style = style
	return logList
}

// SetItemStyle implements LogList.SetItemStyle interface.
func (logList *logListImpl) SetItemStyle(style lipgloss.Style) LogList {
	logList.itemStyle = style
	return logList
}

// NewList creates a new LogList instance.
func NewList() LogList {
	return &logListImpl{}
}

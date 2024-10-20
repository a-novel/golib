package formatters

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/samber/lo"
)

// LogList is a utility LogContent that renders a list of LogContent instances.
type LogList interface {
	LogContent
	// Append adds a new string to the list.
	Append(content string) LogList
	// Nest adds a new LogContent to the list.
	Nest(content LogContent) LogList

	// SetEnumerator sets the list enumerator.
	SetEnumerator(enumerator list.Enumerator) LogList
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
	// ligloss list style.
	style lipgloss.Style
	// ligloss list item style.
	itemStyle lipgloss.Style
}

// RenderConsole implements LogContent.RenderConsole interface.
func (l *logListImpl) RenderConsole() string {
	rendered := list.New(
		lo.Map(l.elements, func(item interface{}, _ int) interface{} {
			if renderer, ok := item.(LogContent); ok {
				return renderer.RenderConsole()
			}
			return item
		})...,
	).
		Enumerator(l.enumerator).
		EnumeratorStyle(l.style).
		ItemStyle(l.itemStyle)

	return rendered.String()
}

// RenderJSON implements LogContent.RenderJSON interface..
func (l *logListImpl) RenderJSON() interface{} {
	return lo.Map(l.elements, func(item interface{}, _ int) interface{} {
		if renderer, ok := item.(LogContent); ok {
			return renderer.RenderJSON()
		}
		return item
	})
}

// Append implements LogList.Append interface.
func (l *logListImpl) Append(content string) LogList {
	l.elements = append(l.elements, NewBase(content))
	return l
}

// Nest implements LogList.Nest interface.
func (l *logListImpl) Nest(content LogContent) LogList {
	l.elements = append(l.elements, content)
	return l
}

// SetEnumerator implements LogList.SetEnumerator interface.
func (l *logListImpl) SetEnumerator(enumerator list.Enumerator) LogList {
	l.enumerator = enumerator
	return l
}

// SetStyle implements LogList.SetStyle interface.
func (l *logListImpl) SetStyle(style lipgloss.Style) LogList {
	l.style = style
	return l
}

// SetItemStyle implements LogList.SetItemStyle interface.
func (l *logListImpl) SetItemStyle(style lipgloss.Style) LogList {
	l.itemStyle = style
	return l
}

// NewList creates a new LogList instance.
func NewList() LogList {
	return &logListImpl{}
}

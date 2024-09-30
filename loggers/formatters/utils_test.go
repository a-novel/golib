package formatters_test

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Ensure tests are termina-agnostic.
func configureTerminal() func() {
	initialColorProfile := termenv.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)

	return func() {
		lipgloss.SetColorProfile(initialColorProfile)
	}
}

// Escape ansi characters for printing debug values.
func escapeAnsi(s string) string {
	s = strings.ReplaceAll(s, "\x1b", "\\x1b")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

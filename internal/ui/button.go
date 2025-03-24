package ui

import (
	"github.com/charmbracelet/lipgloss"
)

type Button struct {
	Label        string
	Focused      bool
	Style        lipgloss.Style
	FocusedStyle lipgloss.Style
}

// View renders the button
func (b Button) View() string {
	if b.Focused {
		return b.FocusedStyle.Render(b.Label)
	}
	return b.Style.Render(b.Label)
}

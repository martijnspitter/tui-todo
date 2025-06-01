package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

const cactus = `
				 /||\
				 ||||
				 ||||
				 |||| /|\
			/|\  |||| |||
			|||  |||| |||
			|||  |||| |||
			|||  |||| d||
			|||  |||||||/
			||b._||||~~'
			\||||||||
			'~~~||||
				||||
				||||
~~~~~~~~~~~~~~~~||||~~~~~~~~~~~~~~
	\/..__..--  . |||| \/  .  ..
		\/         \/ \/    \/
`

// EmptyNothingFoundView creates an empty state view for when no items are found in a search or filter
func EmptyNothingFoundView(translator *i18n.TranslationService, width, height int) string {
	emptyTitle := styling.EmptyStyle.Bold(true).Width(width).Align(lipgloss.Center).Render(translator.T("feedback.nothing_found"))
	cactusText := styling.EmptyStyle.PaddingRight(2).Render(cactus)
	emptyState := lipgloss.JoinVertical(lipgloss.Center, emptyTitle, cactusText)

	return lipgloss.Place(
		width,
		height+2,
		lipgloss.Center,
		lipgloss.Center,
		emptyState,
	)
}

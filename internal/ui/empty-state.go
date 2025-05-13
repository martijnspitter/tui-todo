package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

const trophy = `
    ____________
   '._==_==_==_.'
   .-\:       /-.
  | (|:.      |) |
   '-|:.      |-'
     \::.     /
      '::.  .'
        )  (
      _.'  '._
     '""""""""'
     `

func EmptyStateView(translator *i18n.TranslationService, width, height int) string {
	emptyTitle := styling.EmptyStyle.Bold(true).Width(width).Align(lipgloss.Center).Render(translator.T("feedback.no_todos"))
	emptySubTitle := styling.EmptyStyle.Bold(true).Width(width).Align(lipgloss.Center).Render(translator.T("feedback.mission_accomplished"))
	trophyText := styling.EmptyStyle.PaddingRight(2).Render(trophy)
	emptyState := lipgloss.JoinVertical(lipgloss.Center, emptyTitle, emptySubTitle, trophyText)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		emptyState,
	)
}

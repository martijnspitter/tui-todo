package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type FooterModel struct {
	service    *service.AppService
	tuiService *service.TuiService
	translator *i18n.TranslationService
	width      int
	height     int
	help       tea.Model
}

func NewFooterModel(service *service.AppService, tuiService *service.TuiService, translationService *i18n.TranslationService) *FooterModel {
	help := NewHelpModel(service, tuiService, translationService)
	return &FooterModel{
		service:    service,
		tuiService: tuiService,
		translator: translationService,
		help:       help,
	}
}

func (m *FooterModel) Init() tea.Cmd {
	return nil
}

func (m *FooterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.tuiService.KeyMap.Help):
			m.help.(*HelpModel).ToggleShowAll()
		}
	}
	return m, nil
}

func (m *FooterModel) View() string {
	// Base style for the footer
	footerStyle := lipgloss.NewStyle().
		Background(styling.BackgroundColor).
		Foreground(styling.TextColor).
		Width(m.width)

	content := ""

	// Only show filter options when in specific modes
	if m.tuiService.FilterState.Mode == service.AllFilter ||
		m.tuiService.FilterState.Mode == service.TagFilter {

		// Filter option style
		filterOptionStyle := lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

		// Active filter style
		activeFilterStyle := lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Background(styling.Lavender).
			Foreground(styling.BlackColor)

		// Create option for archived toggle
		var archivedOption string
		if m.tuiService.FilterState.IncludeArchived {
			archivedOption = activeFilterStyle.Render(m.translator.T("footer.show_archived"))
		} else {
			archivedOption = filterOptionStyle.Render(m.translator.T("footer.hide_archived"))
		}

		// Collect all filter options
		filterOptions := []string{archivedOption}

		// Add more filter options as needed

		// Combine options with separator
		separator := lipgloss.NewStyle().
			Foreground(styling.SubtextColor).
			Render(" | ")

		content = archivedOption
		for i := 1; i < len(filterOptions); i++ {
			content = lipgloss.JoinHorizontal(lipgloss.Center, content, separator, filterOptions[i])
		}
	}

	// Join everything
	helpText := m.help.View()
	if m.tuiService.ShowConfirmQuit {
		helpText = styling.SubtextStyle.Render("Really quit? (Press ctrl+c/esc again to quit)")
	}
	statusBar := lipgloss.JoinHorizontal(lipgloss.Left, content)

	return lipgloss.JoinVertical(lipgloss.Center, helpText, footerStyle.Render(statusBar))
}

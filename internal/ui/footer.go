package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type FooterModel struct {
	service    *service.AppService
	tuiService *service.TuiService
	width      int
	height     int
}

func NewFooterModel(service *service.AppService, tuiService *service.TuiService) *FooterModel {
	return &FooterModel{
		service:    service,
		tuiService: tuiService,
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
			archivedOption = activeFilterStyle.Render("[A] Archived")
		} else {
			archivedOption = filterOptionStyle.Render("[A] Archived")
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
	fullContent := lipgloss.JoinHorizontal(lipgloss.Left, content)

	return footerStyle.Render(fullContent)
}

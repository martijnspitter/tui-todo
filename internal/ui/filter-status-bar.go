package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type FilterStatusBar struct {
	tuiService *service.TuiService
	width      int
}

func NewFilterStatusBar(tuiService *service.TuiService) *FilterStatusBar {
	return &FilterStatusBar{
		tuiService: tuiService,
	}
}

func (m *FilterStatusBar) SetWidth(width int) {
	m.width = width
}

func (m *FilterStatusBar) Init() tea.Cmd {
	return nil
}

func (m *FilterStatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

func (m *FilterStatusBar) View() string {
	// Only show filters when in specific modes
	if m.tuiService.FilterState.Mode != service.AllFilter &&
		m.tuiService.FilterState.Mode != service.TagFilter {
		return ""
	}

	// Base style for the status bar
	statusBarStyle := lipgloss.NewStyle().
		Background(styling.BackgroundColor).
		Foreground(styling.TextColor).
		Padding(0, 1).
		MarginTop(1).
		Width(m.width - 2)

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

	// Add more filter options as needed, for example:
	// - Sort by options
	// - Due date filtering
	// - Priority filtering

	// Combine options with separator
	separator := lipgloss.NewStyle().
		Foreground(styling.SubtextColor).
		Render(" | ")

	content := archivedOption
	for i := 1; i < len(filterOptions); i++ {
		content = lipgloss.JoinHorizontal(lipgloss.Center, content, separator, filterOptions[i])
	}

	return statusBarStyle.Render(content)
}

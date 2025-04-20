package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type StatusBar struct {
	service    *service.AppService
	tuiService *service.TuiService
	translator *i18n.TranslationService
	width      int
	height     int
}

func NewStatusBar(service *service.AppService, tuiService *service.TuiService, translator *i18n.TranslationService) *StatusBar {
	return &StatusBar{
		service:    service,
		tuiService: tuiService,
		translator: translator,
	}
}

func (m *StatusBar) Init() tea.Cmd {
	return nil
}

func (m *StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *StatusBar) View() string {
	updateInfo := m.service.GetUpdateInfo()

	content := ""

	// Base style for the status bar
	statusBarStyle := lipgloss.NewStyle().
		Background(styling.BackgroundColor).
		Foreground(styling.TextColor).
		Width(m.width)

	versionStyle := lipgloss.NewStyle().
		Background(styling.BackgroundColor).
		Foreground(styling.TextColor)

	// Filter option style
	filterOptionStyle := lipgloss.NewStyle().
		Background(styling.BackgroundColor).
		Foreground(styling.Lavender).
		PaddingLeft(1).
		PaddingRight(1)

	// Create option for archived toggle
	var archivedOption string
	if m.tuiService.CurrentView == service.AllPane {
		if m.tuiService.FilterState.IncludeArchived {
			archivedOption = filterOptionStyle.Render(m.translator.T("footer.show_archived"))
		} else {
			archivedOption = filterOptionStyle.Render(m.translator.T("footer.hide_archived"))
		}
	}

	// Collect all filter options
	filterOptions := []string{archivedOption}

	if m.tuiService.IsTagFilterActive() {
		tagFilter := filterOptionStyle.
			Render(" 🏷️ " + m.translator.T("filter.by_tag"))
		filterOptions = append(filterOptions, tagFilter)
	}

	if m.tuiService.IsTitleFilterActive() {
		titleFilter := filterOptionStyle.
			Render(" 🔍 " + m.translator.T("filter.by_title"))
		filterOptions = append(filterOptions, titleFilter)
	}

	for i := 1; i < len(filterOptions); i++ {
		content = lipgloss.JoinHorizontal(lipgloss.Center, content, filterOptions[i])
	}

	version := ""
	if updateInfo.HasUpdate {
		version = versionStyle.Render(
			" 🔔",
			m.translator.T("update_available"),
		)
	} else if updateInfo.Version != "" {
		version = versionStyle.Render(
			m.translator.Tf("version", map[string]interface{}{"Version": updateInfo.Version}),
		)
	}

	remainingWidth := m.width - lipgloss.Width(content) - lipgloss.Width(archivedOption) - lipgloss.Width(version) - 4

	spacer := ""
	if remainingWidth > 0 {
		spacer = filterOptionStyle.Render(strings.Repeat(" ", remainingWidth))
	}

	content = lipgloss.JoinHorizontal(lipgloss.Center, content, spacer, archivedOption, version)

	return statusBarStyle.Render(content)
}

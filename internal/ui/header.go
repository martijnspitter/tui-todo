package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
	"github.com/martijnspitter/tui-todo/internal/theme"
)

// HeaderModel represents the header component of the application
type HeaderModel struct {
	tuiService *service.TuiService
	translator *i18n.TranslationService
	width      int
	height     int
}

// NewHeaderModel creates a new header component
func NewHeaderModel(tuiService *service.TuiService, translationService *i18n.TranslationService) *HeaderModel {
	return &HeaderModel{
		tuiService: tuiService,
		translator: translationService,
	}
}

func (m *HeaderModel) Init() tea.Cmd {
	return nil
}

func (m *HeaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *HeaderModel) View() string {
	var leftTabs []string

	isTodaySelected := m.tuiService.CurrentView == service.TodayPane
	todayTab := styling.GetStyledTagWithIndicator(1, m.translator.T("filter.today"), theme.Lavender, isTodaySelected, false, false)
	leftTabs = append(leftTabs, todayTab)

	for status := models.Open; status <= models.Blocked; status++ {
		isSelected := int(m.tuiService.CurrentView) == int(status)+2
		translatedStatus := m.translator.T(status.String())
		tab := styling.GetStyledStatus(translatedStatus, status, isSelected, false, false)
		leftTabs = append(leftTabs, tab)
	}

	leftContent := lipgloss.JoinHorizontal(lipgloss.Center, leftTabs...)

	isAllSelected := m.tuiService.CurrentView == service.AllPane
	allTab := styling.GetStyledTagWithIndicator(6, m.translator.T("filter.all"), theme.Rosewater, isAllSelected, false, false)

	isTagsSelected := m.tuiService.CurrentView == service.TagsPane
	tagsTab := styling.GetStyledTagWithIndicator(7, m.translator.T("filter.tags"), theme.Teal, isTagsSelected, false, false)

	const minGap = 2
	availableWidth := m.width - 2 // -2 for padding
	leftWidth := lipgloss.Width(leftContent)
	rightWidth := lipgloss.Width(allTab) + lipgloss.Width(tagsTab)

	if leftWidth+minGap+rightWidth >= availableWidth {
		return lipgloss.JoinHorizontal(lipgloss.Center, leftContent, allTab, tagsTab)
	}

	spacerWidth := availableWidth - leftWidth - rightWidth
	spacer := strings.Repeat(" ", spacerWidth)

	return lipgloss.JoinHorizontal(lipgloss.Center, leftContent, spacer, allTab, tagsTab)
}

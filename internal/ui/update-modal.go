package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type UpdateModal struct {
	notes      string
	width      int
	height     int
	tuiService *service.TuiService
	translator *i18n.TranslationService
}

func NewUpdateModal(notes string, width, height int, tuiService *service.TuiService, translator *i18n.TranslationService) *UpdateModal {
	return &UpdateModal{
		notes:      notes,
		width:      width,
		height:     height,
		tuiService: tuiService,
		translator: translator,
	}
}

func (m *UpdateModal) Init() tea.Cmd {
	return nil
}

func (m *UpdateModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(
			msg,
			m.tuiService.KeyMap.Quit,
		) {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *UpdateModal) View() string {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width((m.width / 3) * 2).
		BorderForeground(styling.Mauve)

	title := styling.FocusedStyle.Render(m.translator.T("update_required"))
	subtitle := styling.TextStyle.Render(m.translator.T("update_required_subtitle"))
	spacer := ""

	notes := renderReleaseNotes(m.notes)

	content := lipgloss.JoinVertical(lipgloss.Left, title, spacer, subtitle, spacer, notes)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(content),
	)
}

func renderReleaseNotes(markdown string) string {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	rendered, err := r.Render(markdown)
	if err != nil {
		// Fallback to raw markdown if rendering fails
		return markdown
	}

	return rendered
}

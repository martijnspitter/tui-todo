package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
	"github.com/martijnspitter/tui-todo/internal/theme"
)

const Logo = `
 _____         _         _____ _   _ ___
|_   _|__   __| | ___   |_   _| | | |_ _|
  | |/ _ \ / _  |/ _ \    | | | | | || |
  | | (_) | (_| | (_) |   | | | |_| || |
  |_|\___/ \__,_|\___/    |_|  \___/|___|
`

type AboutModal struct {
	width      int
	height     int
	appService *service.AppService
	tuiService *service.TuiService
	translator *i18n.TranslationService
	help       tea.Model
}

func NewAboutModal(width, height int, appService *service.AppService, tuiService *service.TuiService, translator *i18n.TranslationService) *AboutModal {
	help := NewHelpModel(appService, tuiService, translator)

	return &AboutModal{
		width:      width,
		height:     height,
		appService: appService,
		tuiService: tuiService,
		translator: translator,
		help:       help,
	}
}

func (m *AboutModal) Init() tea.Cmd {
	return nil
}

func (m *AboutModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(
			msg,
			m.tuiService.KeyMap.Quit,
		) {
			return m, func() tea.Msg { return modalCloseMsg{reload: false} }
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *AboutModal) View() string {
	updateInfo := m.appService.GetUpdateInfo()

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width((m.width / 3) * 2).
		BorderForeground(theme.Mauve)

	title := styling.FocusedStyle.Render(m.translator.T("about_title"))
	subtitle := styling.TextStyle.Render(m.translator.T("about_subtitle"))
	url := styling.TextStyle.Render(updateInfo.URL)
	spacer := ""
	notesTitle := styling.FocusedStyle.Render(m.translator.Tf("about_notes_title", map[string]interface{}{"Version": updateInfo.Version}))
	notes := styling.RenderMarkdown(updateInfo.Notes)

	help := m.help.View()
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		Logo,
		title,
		spacer,
		subtitle,
		spacer,
		url,
		spacer,
		notesTitle,
		spacer,
		notes,
		spacer,
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(content),
	)
}

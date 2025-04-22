package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/theme"
)

type FooterModel struct {
	service    *service.AppService
	tuiService *service.TuiService
	translator *i18n.TranslationService
	width      int
	height     int
	help       tea.Model
	statusBar  tea.Model
}

func NewFooterModel(service *service.AppService, tuiService *service.TuiService, translationService *i18n.TranslationService) *FooterModel {
	help := NewHelpModel(service, tuiService, translationService)
	statusBar := NewStatusBar(service, tuiService, translationService)
	return &FooterModel{
		service:    service,
		tuiService: tuiService,
		translator: translationService,
		help:       help,
		statusBar:  statusBar,
	}
}

func (m *FooterModel) Init() tea.Cmd {
	return nil
}

func (m *FooterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

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

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	m.statusBar, cmd = m.statusBar.Update(msg)
	cmds = append(cmds, cmd)

	return m, nil
}

func (m *FooterModel) View() string {
	// Join everything
	helpText := m.help.View()
	if m.tuiService.ShowConfirmQuit {
		helpText = lipgloss.NewStyle().Foreground(theme.HelpTextColor).Render("Really quit? (Press ctrl+c/esc again to quit)")
	}

	statusBar := m.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Center, helpText, statusBar)
}

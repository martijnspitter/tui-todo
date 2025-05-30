package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
	"github.com/martijnspitter/tui-todo/internal/theme"
)

type tagState int

const (
	browsingTags tagState = iota
	creatingTag
	deletingTag
)

// TagEditModal allows viewing and managing tags
type TagEditModal struct {
	tag        *models.Tag
	nameInput  textinput.Model
	width      int
	height     int
	appService *service.AppService
	tuiService *service.TuiService
	translator *i18n.TranslationService
	help       tea.Model
}

func NewTagEditModal(tag *models.Tag, width, height int, appService *service.AppService, tuiService *service.TuiService, translationService *i18n.TranslationService) *TagEditModal {
	help := NewHelpModel(appService, tuiService, translationService)

	// Setup tag input for creating new tags
	ti := textinput.New()
	ti.Placeholder = translationService.T("tag.new_placeholder")
	ti.CharLimit = 50

	return &TagEditModal{
		nameInput:  ti,
		width:      width,
		height:     height,
		appService: appService,
		tuiService: tuiService,
		translator: translationService,
		help:       help,
	}
}

func (m *TagEditModal) Init() tea.Cmd {
	return textinput.Blink
}

func (m *TagEditModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.tuiService.KeyMap.Quit):
			// Close modal without saving
			return m, func() tea.Msg { return modalCloseMsg{reload: false} }
		case key.Matches(msg, m.tuiService.KeyMap.AdvanceStatus):
			return m, m.saveChangesCmd()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.nameInput, cmd = m.nameInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *TagEditModal) View() string {
	// Create modal style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(m.width / 2).
		BorderForeground(theme.Mauve)

	content := fmt.Sprintf(
		"%s\n\n%s\n%s\n\n%s",
		styling.FocusedStyle.Render(m.translator.T("tag.create_new")),
		m.nameInput.View(),
		styling.TextStyle.Render(m.translator.T("tag.create_hint")),
		m.help.View(),
	)

	// Center the modal
	positioned := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(content),
	)

	return positioned
}

// ===========================================================================
// Commands
// ===========================================================================
func (m *TagEditModal) saveChangesCmd() tea.Cmd {
	return func() tea.Msg {
		name := m.nameInput.Value()
		if name == "" {
			return TodoErrorMsg{err: fmt.Errorf("tag name cannot be empty")}
		}

		if m.tag != nil {
			// Update existing tag
			m.tag.Name = name
			err := m.appService.UpdateTag(m.tag)
			if err != nil {
				return TodoErrorMsg{err: err}
			}
		} else {
			// Create new tag
			err := m.appService.CreateTag(name)
			if err != nil {
				return TodoErrorMsg{err: err}
			}
		}

		return modalCloseMsg{reload: true}
	}
}

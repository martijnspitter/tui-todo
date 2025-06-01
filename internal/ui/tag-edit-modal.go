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

// TagEditModal allows viewing and managing tags
type TagEditModal struct {
	tag               *models.Tag
	nameInput         textinput.Model
	descInput         textinput.Model
	nameInputSelected bool
	width             int
	height            int
	appService        *service.AppService
	tuiService        *service.TuiService
	translator        *i18n.TranslationService
	help              tea.Model
}

func NewTagEditModal(tag *models.Tag, width, height int, appService *service.AppService, tuiService *service.TuiService, translationService *i18n.TranslationService) *TagEditModal {
	help := NewHelpModel(appService, tuiService, translationService)

	// Setup tag input for creating new tags
	nameInput := textinput.New()
	nameInput.Placeholder = translationService.T("field.name_placeholder")
	nameInput.SetValue(tag.Name)
	nameInput.Focus()
	descInput := textinput.New()
	descInput.SetValue(tag.Description)
	descInput.Placeholder = translationService.T("field.description_placeholder")

	return &TagEditModal{
		tag:               tag,
		nameInput:         nameInput,
		descInput:         descInput,
		nameInputSelected: true,
		width:             width,
		height:            height,
		appService:        appService,
		tuiService:        tuiService,
		translator:        translationService,
		help:              help,
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
		case key.Matches(msg, m.tuiService.KeyMap.Next):
			// Switch focus between name and description inputs
			if m.nameInputSelected {
				m.nameInputSelected = false
				m.nameInput.Blur()
				m.descInput.Focus()
			} else {
				m.nameInputSelected = true
				m.descInput.Blur()
				m.nameInput.Focus()
			}
		case key.Matches(msg, m.tuiService.KeyMap.Prev):
			// Switch focus back to name input
			if !m.nameInputSelected {
				m.nameInputSelected = true
				m.descInput.Blur()
				m.nameInput.Focus()
			} else {
				m.nameInputSelected = false
				m.nameInput.Blur()
				m.descInput.Focus()
			}
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

	if m.nameInputSelected {
		m.nameInput, cmd = m.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.descInput, cmd = m.descInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *TagEditModal) View() string {
	// Create modal style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(m.width / 2).
		BorderForeground(theme.Mauve)

	title := m.translator.T("modal.new_tag")
	if m.tag.ID >= 0 {
		title = m.translator.Tf("modal.edit_tag", map[string]interface{}{"ID": m.tag.ID})
	}
	header := styling.TextStyle.Render(title)

	nameTitle := m.translator.T("field.name")
	descTitle := m.translator.T("field.description")
	if m.nameInputSelected {
		nameTitle = styling.FocusedStyle.Render(nameTitle)
	} else {
		descTitle = styling.FocusedStyle.Render(descTitle)
	}

	content := fmt.Sprintf(
		"%s\n\n%s\n%s\n\n%s\n%s\n\n%s",
		header,
		nameTitle,
		m.nameInput.View(),
		descTitle,
		m.descInput.View(),
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
			return TodoErrorMsg{err: fmt.Errorf(m.translator.T("error.tag_name_empty"))}
		}

		if m.tag.ID >= 0 {
			// Update existing tag
			m.tag.Name = name
			m.tag.Description = m.descInput.Value()
			err := m.appService.UpdateTag(m.tag)
			if err != nil {
				return TodoErrorMsg{err: err}
			}
		} else {
			// Create new tag
			tag := &models.Tag{
				Name:        name,
				Description: m.descInput.Value(),
			}
			err := m.appService.CreateTag(tag)
			if err != nil {
				return TodoErrorMsg{err: err}
			}
		}

		return tagModalCloseMsg{reload: true}
	}
}

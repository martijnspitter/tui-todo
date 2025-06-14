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

type ConfirmDeleteModel struct {
	service      *service.AppService
	tuiService   *service.TuiService
	translator   *i18n.TranslationService
	cancelButton Button
	sendButton   Button
	focused      int
	entityID     int64
	width        int
	height       int
	deleteTag    bool
}

func NewConfirmDeleteModal(appService *service.AppService, tuiService *service.TuiService, translationService *i18n.TranslationService, entityID int64, deleteTag bool) *ConfirmDeleteModel {
	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.SubtextColor)).
		Padding(0, 1)

	focusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.TextColor)).
		Background(theme.SuccessColor).
		Padding(0, 1)

	cancelFocusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.TextColor)).
		Background(theme.ErrorColor).
		Padding(0, 1)

	return &ConfirmDeleteModel{
		service:    appService,
		tuiService: tuiService,
		translator: translationService,
		entityID:   entityID,
		cancelButton: Button{
			Label:        translationService.T("button.cancel"),
			Focused:      true,
			Style:        normalStyle,
			FocusedStyle: cancelFocusedStyle,
		},
		sendButton: Button{
			Label:        translationService.T("button.delete"),
			Style:        normalStyle,
			FocusedStyle: focusedStyle,
		},
		focused:   0, // Start with cancel button focused
		deleteTag: deleteTag,
	}
}

func (m *ConfirmDeleteModel) Init() tea.Cmd {
	return nil
}

func (m *ConfirmDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(
			msg,
			m.tuiService.KeyMap.Next,
			m.tuiService.KeyMap.Prev,
		):
			// Toggle focus between buttons
			if m.focused == 0 {
				m.focused = 1
			} else {
				m.focused = 0
			}
			m.cancelButton.Focused = m.focused == 0
			m.sendButton.Focused = m.focused == 1
			return m, nil

		case key.Matches(msg, m.tuiService.KeyMap.Select):
			// If enter is pressed, trigger the appropriate command based on focus
			if m.focused == 0 {
				return m, CloseModalCmd(false)
			}
			if m.deleteTag {
				// If deleting a tag, call the deleteTagCmd
				return m, m.deleteTagCmd()
			} else {
				// If deleting a todo, call the deleteTodoCmd
				return m, m.deleteTodoCmd()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *ConfirmDeleteModel) View() string {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(m.width / 2).
		BorderForeground(theme.Mauve)

	// Render buttons with appropriate styles
	cancelView := m.cancelButton.View()
	sendView := m.sendButton.View()
	text := m.translator.T("modal.confirm_delete")
	if m.deleteTag {
		text = m.translator.T("modal.confirm_delete_tag")
	}

	title := styling.
		FocusedStyle.
		Width(m.width / 2).
		AlignHorizontal(lipgloss.Center).
		MarginBottom(2).
		Render(text)
	buttons := styling.
		FocusedStyle.
		Width(m.width/2).
		AlignHorizontal(lipgloss.Center).
		Render(cancelView, "  ", sendView)

	content := lipgloss.JoinVertical(lipgloss.Center, title, buttons)

	// Join the buttons with some spacing
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(content),
	)
}

// ===========================================================================
// Messages
// ===========================================================================
type todoDeletedMsg struct{}
type tagDeletedMsg struct{}

// ===========================================================================
// Commands
// ===========================================================================
func (m *ConfirmDeleteModel) deleteTodoCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.service.DeleteTodo(m.entityID)
		if err != nil {
			return TodoErrorMsg{err: err}
		}
		return todoDeletedMsg{}
	}
}

func (m *ConfirmDeleteModel) deleteTagCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.service.DeleteTag(m.entityID)
		if err != nil {
			return TodoErrorMsg{err: err}
		}
		return tagDeletedMsg{}
	}
}

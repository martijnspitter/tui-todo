package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type todoDeletedMsg struct{}

type ConfirmDeleteModel struct {
	service      *service.AppService
	tuiService   *service.TuiService
	cancelButton Button
	sendButton   Button
	focused      int
	todoID       int64
	width        int
	height       int
}

func NewConfirmDeleteModal(appService *service.AppService, tuiService *service.TuiService, todoID int64) *ConfirmDeleteModel {
	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styling.SubtextColor)).
		Padding(0, 1)

	focusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styling.TextColor)).
		Background(styling.SuccessColor).
		Padding(0, 1)

	cancelFocusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styling.TextColor)).
		Background(styling.ErrorColor).
		Padding(0, 1)

	return &ConfirmDeleteModel{
		service:    appService,
		tuiService: tuiService,
		todoID:     todoID,
		cancelButton: Button{
			Label:        "Cancel",
			Focused:      true,
			Style:        normalStyle,
			FocusedStyle: cancelFocusedStyle,
		},
		sendButton: Button{
			Label:        "Delete",
			Style:        normalStyle,
			FocusedStyle: focusedStyle,
		},
		focused: 0, // Start with cancel button focused
	}
}

func (m *ConfirmDeleteModel) deleteTodoCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.service.DeleteTodo(m.todoID)
		if err != nil {
			return todoErrorMsg{err: err}
		}
		return todoDeletedMsg{}
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
			return m, m.deleteTodoCmd()
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
		BorderForeground(styling.Mauve)

	// Render buttons with appropriate styles
	cancelView := m.cancelButton.View()
	sendView := m.sendButton.View()

	title := styling.
		FocusedStyle.
		Width(m.width / 2).
		AlignHorizontal(lipgloss.Center).
		MarginBottom(2).
		Render("Are you sure you want to delete this todo?")
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

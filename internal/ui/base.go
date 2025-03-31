package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/service"
)

type BaseModel struct {
	viewport     viewport.Model
	content      tea.Model
	toastOverlay tea.Model
	ready        bool
}

func NewBaseModel(service *service.AppService, translationService *i18n.TranslationService) *BaseModel {
	todoModel := NewTodoModel(service, translationService)
	toastOverlay := NewToastOverlay(todoModel)
	return &BaseModel{
		toastOverlay: toastOverlay,
		content:      todoModel,
	}
}

func (m *BaseModel) Init() tea.Cmd {
	return nil
}

func (m *BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.content, cmd = m.content.Update(msg)
	cmds = append(cmds, cmd)

	m.toastOverlay, cmd = m.toastOverlay.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Calculate viewport dimensions
		viewportHeight := msg.Height
		viewportWidth := msg.Width

		m.viewport.Width = viewportWidth
		m.viewport.Height = viewportHeight
		m.content.Update(msg)
		m.toastOverlay.Update(msg)

		m.ready = true
	}

	m.viewport.SetContent(m.toastOverlay.View())

	return m, tea.Batch(cmds...)
}

func (m *BaseModel) View() string {
	if !m.ready {
		return "Initializing..."
	}

	return m.viewport.View()
}

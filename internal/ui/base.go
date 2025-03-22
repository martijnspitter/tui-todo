package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/martijnspitter/tui-todo/internal/service"
)

type BaseModel struct {
	viewport viewport.Model
	content  tea.Model
}

func NewBaseModel(service *service.AppService) *BaseModel {
	todoModel := NewTodoModel(service)
	toastOverlay := NewToastOverlay(todoModel)
	return &BaseModel{
		content: toastOverlay,
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

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.content, cmd = m.content.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Calculate viewport dimensions
		viewportHeight := msg.Height
		viewportWidth := msg.Width

		m.viewport.Width = viewportWidth
		m.viewport.Height = viewportHeight

		// After updating dimensions, refresh content
		m.viewport.SetContent(m.content.View())
	}

	return m, tea.Batch(cmds...)
}

func (m *BaseModel) View() string {
	return m.viewport.View()
}

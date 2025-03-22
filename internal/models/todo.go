package models

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/martijnspitter/tui-todo/internal/service"
)

type TodoModel struct {
	tuiService *service.TuiService
	quiting    bool
}

func NewTodoModel() *TodoModel {
	tuiService := service.NewTuiService()
	return &TodoModel{
		tuiService: tuiService,
	}
}

func (m *TodoModel) Init() tea.Cmd {
	return nil
}

func (m *TodoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(
			msg,
			m.tuiService.KeyMap.Quit,
		):
			m.quiting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *TodoModel) View() string {
	if m.quiting {
		return ""
	}

	return "Hello World"
}

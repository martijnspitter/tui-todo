package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type CreateTodoMsg struct {
	Priority models.Priority
	Title    string
}

func CreateTodoCmd(title string, priority models.Priority) tea.Cmd {
	return func() tea.Msg {
		return CreateTodoMsg{
			Priority: priority,
			Title:    title,
		}
	}
}

type selectedInput int

const (
	low selectedInput = iota
	medium
	high
	title
)

type FooterModel struct {
	service       *service.AppService
	tuiService    *service.TuiService
	textInput     textinput.Model
	selectedInput selectedInput
	width         int
	height        int
	priority      models.Priority
}

func NewFooterModel(service *service.AppService, tuiService *service.TuiService) *FooterModel {
	ti := textinput.New()
	ti.Placeholder = "Create new todo..."

	return &FooterModel{
		textInput:  ti,
		priority:   models.Medium,
		service:    service,
		tuiService: tuiService,
	}
}

func (m *FooterModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *FooterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.tuiService.KeyMap.Quit):
			m.tuiService.SelectedPane = service.SelectedPane(msg.String()[0])
		case key.Matches(msg, m.tuiService.KeyMap.Next):
			if m.selectedInput == title {
				m.SetSelectedPane(low)
			} else {
				m.SetSelectedPane(m.selectedInput + 1)
			}
		case key.Matches(msg, m.tuiService.KeyMap.Prev):
			if m.selectedInput == low {
				m.SetSelectedPane(title)
			} else {
				m.SetSelectedPane(m.selectedInput - 1)
			}
		case key.Matches(msg, m.tuiService.KeyMap.Select):
			if m.selectedInput != title {
				m.priority = models.Priority(m.selectedInput)
				m.SetSelectedPane(title)
			} else {
				cmds = append(cmds, CreateTodoCmd(m.textInput.Value(), m.priority))
				m.textInput.SetValue("")
			}

		}
		if m.selectedInput == title {
			// Update the text input
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	}

	return m, tea.Batch(cmds...)
}

func (m *FooterModel) View() string {
	var priorityTabs []string
	for p := models.Priority(0); p < 3; p++ {
		selected := p == m.priority && m.tuiService.SelectedPane == service.New
		hovered := m.selectedInput == selectedInput(int(p)) && m.tuiService.SelectedPane == service.New

		priorityTab := styling.GetStyledPriority(p, selected, hovered)

		priorityTabs = append(priorityTabs, priorityTab)
	}

	prioritySection := lipgloss.JoinHorizontal(lipgloss.Center, priorityTabs...)

	// Combine the new todo tab, priority tabs, and input field
	inputWidth := m.width - lipgloss.Width(prioritySection) - 4
	inputStyle := lipgloss.NewStyle().Width(inputWidth)

	formattedInput := inputStyle.Render(m.textInput.View())

	inputLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		prioritySection,
		formattedInput,
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().MarginTop(1).Render(inputLine))
}

// SetSelectedPane updates the selected pane
func (m *FooterModel) SetSelectedPane(pane selectedInput) {
	m.selectedInput = pane
	if pane == title {
		m.textInput.Focus()
	} else {
		m.textInput.Blur()
	}
}

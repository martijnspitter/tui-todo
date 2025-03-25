package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type editState int

const (
	editingTitle editState = iota
	editingDescription
	editingTags
	editingDueDate
	editingPriorityLow
	editingPriorityMedium
	editingPriorityHigh
)

// TodoEditModal allows viewing and editing todo details
type TodoEditModal struct {
	todo         *models.Todo
	titleInput   textinput.Model
	descInput    textarea.Model
	tagsInput    textinput.Model
	dueDateInput textinput.Model
	priority     models.Priority
	editState    editState
	width        int
	height       int
	appService   *service.AppService
	tuiService   *service.TuiService
}

func NewTodoEditModal(todo *models.Todo, width, height int, appService *service.AppService, tuiService *service.TuiService) *TodoEditModal {
	ti := textinput.New()
	ti.SetValue(todo.Title)
	ti.Focus()

	desc := textarea.New()
	desc.SetValue(todo.Description)
	desc.ShowLineNumbers = true

	tagsInput := textinput.New()
	tagsInput.SetValue(strings.Join(todo.Tags, ", "))

	dueDateInput := textinput.New()
	dueDateInput.Placeholder = "YYYY-MM-DD HH:MM (e.g. 2023-12-31 15:30)"
	if todo.DueDate != nil {
		dueDateInput.SetValue(todo.DueDate.Format("2006-01-02 15:04"))
	}

	return &TodoEditModal{
		todo:         todo,
		titleInput:   ti,
		descInput:    desc,
		tagsInput:    tagsInput,
		dueDateInput: dueDateInput,
		priority:     todo.Priority,
		width:        width,
		height:       height,
		editState:    editingTitle,
		appService:   appService,
		tuiService:   tuiService,
	}
}

func (m *TodoEditModal) Init() tea.Cmd {
	return textinput.Blink
}

func (m *TodoEditModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(
			msg,
			m.tuiService.KeyMap.Quit,
		):
			// Close modal without saving
			return m, func() tea.Msg { return modalCloseMsg{reload: false} }

		case key.Matches(msg, m.tuiService.KeyMap.Next):
			// Cycle through edit states
			m.goForward()
		case key.Matches(msg, m.tuiService.KeyMap.Prev):
			// Cycle through edit states
			m.goBack()

		case key.Matches(msg, m.tuiService.KeyMap.Select):
			if m.editState == editingPriorityLow {
				m.priority = models.Low
			} else if m.editState == editingPriorityMedium {
				m.priority = models.Medium
			} else if m.editState == editingPriorityHigh {
				m.priority = models.High
			}

		case key.Matches(msg, m.tuiService.KeyMap.AdvanceStatus):
			return m, m.saveChangesCmd()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Update active input based on state
	var cmd tea.Cmd
	switch m.editState {
	case editingTitle:
		m.titleInput, cmd = m.titleInput.Update(msg)
		cmds = append(cmds, cmd)
	case editingDescription:
		m.descInput, cmd = m.descInput.Update(msg)
		cmds = append(cmds, cmd)
	case editingTags:
		m.tagsInput, cmd = m.tagsInput.Update(msg)
		cmds = append(cmds, cmd)
	case editingDueDate:
		m.dueDateInput, cmd = m.dueDateInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *TodoEditModal) View() string {
	// Create modal style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(m.width / 2).
		BorderForeground(styling.Mauve)

	status := styling.GetStyledStatus(m.todo.Status, true, true)

	// Priority display
	var priorityTabs []string
	for p := models.Priority(0); p < 3; p++ {
		selected := p == m.priority
		hovered := m.editState == editState(int(p)+4)

		priorityTab := styling.GetStyledPriority(p, selected, hovered)

		priorityTabs = append(priorityTabs, priorityTab)
	}

	prioritySection := lipgloss.JoinHorizontal(lipgloss.Center, priorityTabs...)

	// Title field
	titleField := "Title"
	if m.editState == editingTitle {
		titleField = styling.FocusedStyle.Render(titleField)
	}
	title := fmt.Sprintf("%s\n%s", titleField, m.titleInput.View())

	// Description field
	descField := "Description"
	if m.editState == editingDescription {
		descField = styling.FocusedStyle.Render(descField)
	}
	m.descInput.SetWidth((m.width / 2) - 4)
	description := fmt.Sprintf("%s\n%s", descField, m.descInput.View())

	// Tags field
	tagsField := "Tags (comma separated)"
	if m.editState == editingTags {
		tagsField = styling.FocusedStyle.Render(tagsField)
	}
	tags := fmt.Sprintf("%s\n%s", tagsField, m.tagsInput.View())

	// Priority field
	priorityField := "Priority"
	if m.editState == editingPriorityLow || m.editState == editingPriorityMedium || m.editState == editingPriorityHigh {
		priorityField = styling.FocusedStyle.Render(priorityField)
	}

	// Due Date field
	dueDateField := "Due Date (YYYY-MM-DD HH:MM or empty to clear)"
	if m.editState == editingDueDate {
		dueDateField = styling.FocusedStyle.Render(dueDateField)
	}
	dueDate := fmt.Sprintf("%s\n%s", dueDateField, m.dueDateInput.View())

	// Combine all content
	content := fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s",
		styling.TextStyle.Render(fmt.Sprintf("Editing Todo #%d", m.todo.ID)),
		status,
		title,
		description,
		tags,
		dueDate,
		fmt.Sprintf("%s\n%s", priorityField, prioritySection),
		styling.SubtextStyle.Render("ctrl+s: save  tab: next field  esc: cancel"),
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

func (m *TodoEditModal) goForward() {
	switch m.editState {
	case editingTitle:
		m.titleInput.Blur()
		m.descInput.Focus()
		m.editState = editingDescription
	case editingDescription:
		m.descInput.Blur()
		m.tagsInput.Focus()
		m.editState = editingTags
	case editingTags:
		m.tagsInput.Blur()
		m.dueDateInput.Focus()
		m.editState = editingDueDate
	case editingDueDate:
		m.dueDateInput.Blur()
		m.editState = editingPriorityLow
	case editingPriorityLow:
		m.editState = editingPriorityMedium
	case editingPriorityMedium:
		m.editState = editingPriorityHigh
	case editingPriorityHigh:
		m.titleInput.Focus()
		m.editState = editingTitle
	}
}

func (m *TodoEditModal) goBack() {
	switch m.editState {
	case editingTitle:
		m.titleInput.Blur()
		m.editState = editingPriorityHigh
	case editingDescription:
		m.descInput.Blur()
		m.titleInput.Focus()
		m.editState = editingTitle
	case editingTags:
		m.descInput.Focus()
		m.tagsInput.Blur()
		m.editState = editingDescription
	case editingDueDate:
		m.tagsInput.Focus()
		m.dueDateInput.Blur()
		m.editState = editingTags
	case editingPriorityLow:
		m.dueDateInput.Focus()
		m.editState = editingDueDate
	case editingPriorityMedium:
		m.editState = editingPriorityLow
	case editingPriorityHigh:
		m.editState = editingPriorityMedium
	}
}

func (m *TodoEditModal) saveChangesCmd() tea.Cmd {
	return func() tea.Msg {
		// Update todo with new values
		m.todo.Title = m.titleInput.Value()
		m.todo.Description = m.descInput.Value()
		m.todo.Priority = m.priority

		dueDateStr := strings.TrimSpace(m.dueDateInput.Value())
		if dueDateStr == "" {
			// Clear due date if empty
			m.todo.DueDate = nil
		} else {
			// Try to parse the date string
			dueDate, err := time.Parse("2006-01-02 15:04", dueDateStr)
			if err != nil {
				return todoErrorMsg{err: fmt.Errorf("invalid due date format: %w", err)}
			}
			m.todo.DueDate = &dueDate
		}

		// Handle tags (split by comma and trim)
		rawTags := strings.Split(m.tagsInput.Value(), ",")
		tags := make([]string, 0, len(rawTags))
		for _, tag := range rawTags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tags = append(tags, tag)
			}
		}

		err := m.appService.SaveTodo(m.todo, tags)
		if err != nil {
			return todoErrorMsg{err: err}
		}

		// Close modal and reload todos
		return modalCloseMsg{reload: true}
	}
}

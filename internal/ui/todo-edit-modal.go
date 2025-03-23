package ui

import (
	"fmt"
	"strings"

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
	editingPriority
	editingTags
)

// TodoEditModal allows viewing and editing todo details
type TodoEditModal struct {
	todo       *models.Todo
	titleInput textinput.Model
	descInput  textarea.Model
	tagsInput  textinput.Model
	priority   models.Priority
	editState  editState
	width      int
	height     int
	appService *service.AppService
}

func NewTodoEditModal(todo *models.Todo, width, height int, appService *service.AppService) *TodoEditModal {
	ti := textinput.New()
	ti.SetValue(todo.Title)
	ti.Focus()

	desc := textarea.New()
	desc.SetValue(todo.Description)

	tagsInput := textinput.New()
	tagsInput.SetValue(strings.Join(todo.Tags, ", "))

	return &TodoEditModal{
		todo:       todo,
		titleInput: ti,
		descInput:  desc,
		tagsInput:  tagsInput,
		priority:   todo.Priority,
		width:      width,
		height:     height,
		appService: appService,
		editState:  editingTitle,
	}
}

func (m *TodoEditModal) Init() tea.Cmd {
	return textinput.Blink
}

func (m *TodoEditModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Close modal without saving
			return m, func() tea.Msg { return modalCloseMsg{reload: false} }

		case "tab":
			// Cycle through edit states
			m.toggleFocus()

		case "enter":
			if m.editState == editingPriority {
				// Cycle priority
				m.priority = (m.priority + 1) % 3
			} else if m.editState == editingTitle && m.titleInput.Value() == "" {
				// Can't have empty title
				return m, func() tea.Msg {
					return todoErrorMsg{err: fmt.Errorf("title cannot be empty")}
				}
			}

		case "ctrl+s":
			// Save changes
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
	}

	return m, tea.Batch(cmds...)
}

func (m *TodoEditModal) View() string {
	// Create modal style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(m.width / 2).
		BorderForeground(styling.BrandColor)

	// Status pill
	statusStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Bold(true).
		Foreground(lipgloss.Color("#000"))

	var statusColor lipgloss.Color
	switch m.todo.Status {
	case models.Open:
		statusColor = lipgloss.Color("#AED6F1") // Light blue
	case models.Doing:
		statusColor = lipgloss.Color("#F9E79F") // Light yellow
	case models.Done:
		statusColor = lipgloss.Color("#ABEBC6") // Light green
	case models.Archived:
		statusColor = lipgloss.Color("#D4EFDF") // Very light green
	}

	statusPill := statusStyle.Copy().Background(statusColor).Render(m.todo.Status.String())

	// Priority display
	priorityDisplay := "Priority: "
	for p := models.Priority(0); p < 3; p++ {
		if p == m.priority {
			priorityDisplay += styling.FocusedStyle.Render(p.String() + " ")
		} else {
			priorityDisplay += p.String() + " "
		}
	}

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
	description := fmt.Sprintf("%s\n%s", descField, m.descInput.View())

	// Tags field
	tagsField := "Tags (comma separated)"
	if m.editState == editingTags {
		tagsField = styling.FocusedStyle.Render(tagsField)
	}
	tags := fmt.Sprintf("%s\n%s", tagsField, m.tagsInput.View())

	// Priority field
	priorityField := "Priority"
	if m.editState == editingPriority {
		priorityField = styling.FocusedStyle.Render(priorityField)
	}

	// Combine all content
	content := fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s",
		fmt.Sprintf("Editing Todo #%d", m.todo.ID),
		statusPill,
		title,
		description,
		tags,
		fmt.Sprintf("%s\n%s", priorityField, priorityDisplay),
		"ctrl+s: save  tab: next field  esc: cancel",
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

func (m *TodoEditModal) toggleFocus() {
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
		m.editState = editingPriority
	case editingPriority:
		m.titleInput.Focus()
		m.editState = editingTitle
	}
}

func (m *TodoEditModal) saveChangesCmd() tea.Cmd {
	return func() tea.Msg {
		// Update todo with new values
		m.todo.Title = m.titleInput.Value()
		m.todo.Description = m.descInput.Value()
		m.todo.Priority = m.priority

		// Handle tags (split by comma and trim)
		rawTags := strings.Split(m.tagsInput.Value(), ",")
		tags := make([]string, 0, len(rawTags))
		for _, tag := range rawTags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tags = append(tags, tag)
			}
		}

		// Update the todo
		err := m.appService.UpdateTodo(m.todo)
		if err != nil {
			return todoErrorMsg{err: err}
		}

		// Handle tag updates (remove all, then add new)
		for _, tag := range tags {
			err := m.appService.AddTagToTodo(m.todo.ID, tag)
			if err != nil {
				return todoErrorMsg{err: err}
			}
		}

		// Close modal and reload todos
		return modalCloseMsg{reload: true}
	}
}

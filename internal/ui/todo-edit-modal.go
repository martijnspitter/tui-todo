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
	"github.com/charmbracelet/log"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
	"github.com/martijnspitter/tui-todo/internal/theme"
	"slices"
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
	editingPriorityMajor
	editingPriorityCritical
	editingStatusOpen
	editingStatusDoing
	editingStatusDone
	editingStatusBlocked
)

// TodoEditModal allows viewing and editing todo details
type TodoEditModal struct {
	todo         *models.Todo
	titleInput   textinput.Model
	descInput    textarea.Model
	tagsInput    *TagSelector
	dueDateInput textinput.Model
	priority     models.Priority
	status       models.Status
	editState    editState
	width        int
	height       int
	appService   *service.AppService
	tuiService   *service.TuiService
	translator   *i18n.TranslationService
	help         tea.Model
}

func NewTodoEditModal(todo *models.Todo, width, height int, appService *service.AppService, tuiService *service.TuiService, translationService *i18n.TranslationService) *TodoEditModal {
	help := NewHelpModel(appService, tuiService, translationService)

	ti := textinput.New()
	ti.SetValue(todo.Title)
	ti.Focus()

	desc := textarea.New()
	desc.SetValue(todo.Description)
	desc.ShowLineNumbers = true

	allTags, err := appService.GetAllTags()
	if err != nil {
		log.Error("Failed to load tags", "error", err)
		allTags = []*models.Tag{}
	}

	// Create tag selector with selected tags
	tagSelector := NewTagSelector(todo.Tags, allTags, tuiService, translationService)

	dueDateInput := textinput.New()
	dueDateInput.Placeholder = "YYYY-MM-DD HH:MM (e.g. 2023-12-31 15:30)"
	if todo.DueDate != nil {
		dueDateInput.SetValue(todo.DueDate.Format("2006-01-02 15:04"))
	}

	return &TodoEditModal{
		todo:         todo,
		titleInput:   ti,
		descInput:    desc,
		tagsInput:    tagSelector,
		dueDateInput: dueDateInput,
		priority:     todo.Priority,
		status:       todo.Status,
		width:        width,
		height:       height,
		editState:    editingTitle,
		appService:   appService,
		tuiService:   tuiService,
		translator:   translationService,
		help:         help,
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
			if m.editState != editingTags {
				// Cycle through edit states
				m.goForward()
			}
		case key.Matches(msg, m.tuiService.KeyMap.Prev):
			if m.editState != editingTags {
				// Cycle through edit states
				m.goBack()
			}
		case key.Matches(msg, m.tuiService.KeyMap.Select):
			switch m.editState {
			case editingPriorityLow:
				m.priority = models.Low
			case editingPriorityMedium:
				m.priority = models.Medium
			case editingPriorityHigh:
				m.priority = models.High
			case editingPriorityMajor:
				m.priority = models.Major
			case editingPriorityCritical:
				m.priority = models.Critical
			case editingStatusOpen:
				m.status = models.Open
			case editingStatusDoing:
				m.status = models.Doing
			case editingStatusDone:
				m.status = models.Done
			case editingStatusBlocked:
				m.status = models.Blocked
			}

		case key.Matches(msg, m.tuiService.KeyMap.AdvanceStatus):
			return m, m.saveChangesCmd()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case GoBackMsg:
		m.goBack()
	case GoForwardMsg:
		m.goForward()
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
		BorderForeground(theme.Mauve)

	// Priority display
	var priorityTabs []string
	for p := models.Priority(0); p <= models.Critical; p++ {
		selected := p == m.priority
		hovered := m.editState == editState(int(p)+4)

		translatedPriority := m.translator.T(p.String())
		priorityTab := styling.GetStyledPriority(translatedPriority, p, selected, hovered)

		priorityTabs = append(priorityTabs, priorityTab)
	}

	prioritySection := lipgloss.JoinHorizontal(lipgloss.Center, priorityTabs...)

	var statusTabs []string
	for status := models.Open; status <= models.Blocked; status++ {
		selected := status == m.status
		hovered := m.editState == editState(int(status)+9)
		prefix := ""
		spacer := " "
		if !selected && !hovered {
			prefix = " "
			spacer = ""
		}
		translatedStatus := m.translator.T(status.String())
		tab := styling.GetStyledStatus(prefix+translatedStatus, status, selected, true, hovered)
		statusTabs = append(statusTabs, tab+spacer)
	}

	statusSection := lipgloss.JoinHorizontal(lipgloss.Center, statusTabs...)

	// Title field
	titleField := m.translator.T("field.title")
	if m.editState == editingTitle {
		titleField = styling.FocusedStyle.Render(titleField)
	}
	title := fmt.Sprintf("%s\n%s", titleField, m.titleInput.View())

	// Description field
	descField := m.translator.T("field.description")
	if m.editState == editingDescription {
		descField = styling.FocusedStyle.Render(descField)
	}
	m.descInput.SetWidth((m.width / 2) - 4)
	description := fmt.Sprintf("%s\n%s", descField, m.descInput.View())

	// Tags field
	tagsField := m.translator.T("field.tags")
	if m.editState == editingTags {
		tagsField = styling.FocusedStyle.Render(tagsField)
	}
	tags := fmt.Sprintf("%s\n%s", tagsField, m.tagsInput.View())

	// Priority header
	priorityHeader := m.translator.T("field.priority")
	if m.editState == editingPriorityLow || m.editState == editingPriorityMedium || m.editState == editingPriorityHigh {
		priorityHeader = styling.FocusedStyle.Render(priorityHeader)
	}

	// Status header
	statusHeader := m.translator.T("field.status")
	if m.editState == editingStatusOpen || m.editState == editingStatusDoing || m.editState == editingStatusDone || m.editState == editingStatusBlocked {
		statusHeader = styling.FocusedStyle.Render(statusHeader)
	}

	// Due Date field
	dueDateField := m.translator.T("field.due_date")
	if m.editState == editingDueDate {
		dueDateField = styling.FocusedStyle.Render(dueDateField)
	}
	dueDate := fmt.Sprintf("%s\n%s", dueDateField, m.dueDateInput.View())

	updatedAtHeader := ""
	updatedAt := ""
	if m.todo.ID >= 0 {
		updatedAtHeader = m.translator.T("field.updated_at")
		translatedUpdatedAt := m.translator.Tf("ui.updated", map[string]interface{}{"Time": m.todo.UpdatedAt.Format(time.Stamp)})
		updatedAtField := styling.GetStyledUpdatedAt(translatedUpdatedAt)
		updatedAt = fmt.Sprintf("%s\n%s", updatedAtHeader, updatedAtField)
	}

	header := m.translator.T("modal.new_todo")
	if m.todo.ID >= 0 {
		text := m.translator.Tf("modal.edit_todo", map[string]interface{}{"ID": m.todo.ID})
		timeSpendText := m.translator.Tf("ui.time_spent", map[string]interface{}{"Time": m.todo.FormatTimeSpent()})
		timeSpend := styling.GetTimeSpend(timeSpendText)

		remainder := m.width/2 - lipgloss.Width(text) - lipgloss.Width(timeSpend) + 8

		header = lipgloss.JoinHorizontal(lipgloss.Left, styling.TextStyle.Width(remainder).Align(lipgloss.Left).Render(text), timeSpend)
	}

	help := m.help.View()

	// Combine all content
	content := fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s",
		styling.TextStyle.Render(header),
		title,
		description,
		tags,
		dueDate,
		fmt.Sprintf("%s\n%s", priorityHeader, prioritySection),
		fmt.Sprintf("%s\n%s", statusHeader, statusSection),
		updatedAt,
		help,
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
// Helpers
// ===========================================================================
func (m *TodoEditModal) goForward() {
	switch m.editState {
	case editingTitle:
		m.titleInput.Blur()
		m.descInput.Focus()
	case editingDescription:
		m.descInput.Blur()
		m.tagsInput.Focus()
	case editingTags:
		m.tagsInput.Blur()
		m.dueDateInput.Focus()
	case editingDueDate:
		m.dueDateInput.Blur()
	case editingStatusBlocked:
		m.titleInput.Focus()
	}

	if m.editState == editingStatusBlocked {
		m.editState = 0
	} else {
		m.editState = m.editState + 1
	}
}

func (m *TodoEditModal) goBack() {
	switch m.editState {
	case editingTitle:
		m.titleInput.Blur()
	case editingDescription:
		m.descInput.Blur()
		m.titleInput.Focus()
	case editingTags:
		m.descInput.Focus()
	case editingDueDate:
		m.tagsInput.Focus()
		m.dueDateInput.Blur()
	case editingPriorityLow:
		m.dueDateInput.Focus()
	}

	if m.editState == editingTitle {
		m.editState = editingStatusBlocked
	} else {
		m.editState = m.editState - 1
	}
}

// ===========================================================================
// Messages
// ===========================================================================
type GoBackMsg struct{}
type GoForwardMsg struct{}

// ===========================================================================
// Commands
// ===========================================================================
func (m *TodoEditModal) saveChangesCmd() tea.Cmd {
	return func() tea.Msg {
		// Update todo with new values
		m.todo.Title = m.titleInput.Value()
		m.todo.Description = m.descInput.Value()
		m.todo.Priority = m.priority
		m.todo.Status = m.status

		dueDateStr := strings.TrimSpace(m.dueDateInput.Value())
		if dueDateStr == "" {
			// Clear due date if empty
			m.todo.DueDate = nil
		} else {
			// Try to parse the date string
			dueDate, err := time.Parse("2006-01-02 15:04", dueDateStr)
			if err != nil {
				log.Error("Invalid due date", err)
				return TodoErrorMsg{err: fmt.Errorf(m.translator.T("error.due_date_invalid"))}
			}
			m.todo.DueDate = &dueDate
		}

		tags := m.tagsInput.SelectedTags()

		err := m.appService.SaveTodo(m.todo, tags)
		if err != nil {
			return TodoErrorMsg{err: err}
		}

		// Close modal and reload todos
		return modalCloseMsg{reload: true}
	}
}

func GoToPreviousEditState() tea.Cmd {
	return func() tea.Msg {
		return GoBackMsg{}
	}
}

func GoToNextEditState() tea.Cmd {
	return func() tea.Msg {
		return GoForwardMsg{}
	}
}

// ===========================================================================
// Tag Selector
// ===========================================================================
type TagSelector struct {
	selectedTags  []string
	availableTags []*models.Tag
	focused       bool
	cursor        int
	tuiService    *service.TuiService
	translator    *i18n.TranslationService
}

func NewTagSelector(selectedTags []string, availableTags []*models.Tag, tuiService *service.TuiService, translator *i18n.TranslationService) *TagSelector {
	return &TagSelector{
		selectedTags:  selectedTags,
		availableTags: availableTags,
		cursor:        -1,
		tuiService:    tuiService,
		translator:    translator,
	}
}

func (ts *TagSelector) Focus() {
	ts.focused = true
}

func (ts *TagSelector) Blur() {
	ts.focused = false
}

func (ts *TagSelector) Update(msg tea.Msg) (*TagSelector, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, ts.tuiService.KeyMap.Prev):
			if ts.cursor > 0 {
				ts.cursor--
			} else {
				ts.cursor = -1
				return ts, GoToPreviousEditState()
			}
		case key.Matches(msg, ts.tuiService.KeyMap.Next):
			if ts.cursor < len(ts.availableTags)-1 {
				ts.cursor++
			} else {
				ts.cursor = len(ts.availableTags)
				return ts, GoToNextEditState()
			}
		case key.Matches(msg, ts.tuiService.KeyMap.Select):
			if ts.cursor < 0 || ts.cursor >= len(ts.availableTags) {
				return ts, nil
			}
			// Toggle selected status of the current tag
			currentTag := ts.availableTags[ts.cursor]
			if ts.IsSelected(currentTag.Name) {
				// Remove tag
				for i, tag := range ts.selectedTags {
					if tag == currentTag.Name {
						ts.selectedTags = slices.Delete(ts.selectedTags, i, i+1)
						break
					}
				}
			} else {
				// Add tag
				ts.selectedTags = append(ts.selectedTags, currentTag.Name)
			}
		}
	}

	return ts, nil
}

func (ts *TagSelector) IsSelected(tag string) bool {
	return slices.Contains(ts.selectedTags, tag)
}

func (ts *TagSelector) View() string {

	if len(ts.availableTags) == 0 {
		return ts.translator.T("feedback.no_tags_available")
	}

	var sb strings.Builder
	sb.WriteString(ts.translator.T("field.select_tags") + "\n\n")

	for i, tag := range ts.availableTags {
		cursor := " "
		if ts.focused && ts.cursor == i {
			cursor = ">"
		}

		checked := "[ ]"
		if ts.IsSelected(tag.Name) {
			checked = "[✓]"
		}

		sb.WriteString(fmt.Sprintf("%s %s %s\n", cursor, checked, styling.GetStyledTag(tag.Name)))
	}

	return sb.String()
}

func (ts *TagSelector) SelectedTags() []string {
	return ts.selectedTags
}

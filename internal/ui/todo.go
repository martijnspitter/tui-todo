package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type TodoItem struct {
	todo *models.Todo
}

func (i TodoItem) Title() string {
	priorityMarker := map[models.Priority]string{
		models.Low:    "⬇️",
		models.Medium: "➡️",
		models.High:   "⬆️",
	}[i.todo.Priority]

	return fmt.Sprintf("%s %s", priorityMarker, i.todo.Title)
}

func (i TodoItem) Description() string {
	var tags string
	if len(i.todo.Tags) > 0 {
		tags = fmt.Sprintf(" [%s]", strings.Join(i.todo.Tags, ", "))
	}

	var dueDate string
	if i.todo.DueDate != nil {
		dueDate = fmt.Sprintf(" (Due: %s)", i.todo.DueDate.Format("Jan 2"))
	}

	return fmt.Sprintf("%s%s%s",
		truncateString(i.todo.Description, 50),
		tags,
		dueDate)
}

func (i TodoItem) FilterValue() string {
	return i.todo.Title
}

type TodoModel struct {
	service        *service.AppService
	tuiService     *service.TuiService
	list           list.Model
	textInput      textinput.Model
	currentFilter  models.Status
	width          int
	height         int
	quitting       bool
	showingModal   bool
	modalComponent tea.Model
	err            error
}

func NewTodoModel(appService *service.AppService) *TodoModel {
	tuiService := service.NewTuiService()

	// Setup input
	ti := textinput.New()
	ti.Placeholder = "Create new todo..."
	ti.Focus()

	// Setup list
	delegate := list.NewDefaultDelegate()
	todoList := list.New([]list.Item{}, delegate, 0, 0)
	todoList.Title = ""
	todoList.SetShowHelp(false)

	// Create model
	m := &TodoModel{
		service:       appService,
		tuiService:    tuiService,
		textInput:     ti,
		list:          todoList,
		currentFilter: models.Doing, // Default to showing "doing" todos
	}

	// Set custom keys
	m.setupKeyMap()

	return m
}

func (m *TodoModel) setupKeyMap() {
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("n"),
				key.WithHelp("n", "new todo"),
			),
			key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("e", "edit todo"),
			),
			key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "delete todo"),
			),
			key.NewBinding(
				key.WithKeys("space"),
				key.WithHelp("space", "advance status"),
			),
			key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "change filter"),
			),
		}
	}
}

func (m *TodoModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadTodos(),
		textinput.Blink,
	)
}

func (m *TodoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showingModal {
			// Handle modal
			var cmd tea.Cmd
			m.modalComponent, cmd = m.modalComponent.Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(
			msg,
			m.tuiService.KeyMap.Quit,
		):
			m.quitting = true
			return m, tea.Quit

		case msg.String() == "tab":
			// Cycle through filters
			m.currentFilter = (m.currentFilter + 1) % 4
			m.list.Title = fmt.Sprintf("Todos - %s", m.currentFilter.String())
			return m, m.loadTodos()

		case msg.String() == "enter":
			// Create new todo when pressing enter in the input
			if m.textInput.Value() != "" {
				cmd := m.createTodoCmd(m.textInput.Value(), "")
				m.textInput.Reset()
				return m, tea.Batch(cmd, m.loadTodos())
			}

		case msg.String() == "space":
			// Advance todo status
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(TodoItem)
				return m, m.advanceTodoStatusCmd(item.todo.ID)
			}

		case msg.String() == "e":
			// Edit selected todo
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(TodoItem)
				return m, m.showEditModalCmd(item.todo)
			}

		case msg.String() == "d":
			// Delete selected todo
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(TodoItem)
				return m, m.deleteTodoCmd(item.todo.ID)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 3 // Title + top border
		footerHeight := 3 // Input + bottom padding

		m.list.SetSize(msg.Width, msg.Height-headerHeight-footerHeight)

		// Pass size to modal if active
		if m.showingModal && m.modalComponent != nil {
			var cmd tea.Cmd
			m.modalComponent, cmd = m.modalComponent.Update(msg)
			cmds = append(cmds, cmd)
		}

	case todosLoadedMsg:
		items := make([]list.Item, len(msg.todos))
		for i, todo := range msg.todos {
			items[i] = TodoItem{todo: todo}
		}
		cmd := m.list.SetItems(items)
		cmds = append(cmds, cmd)

	case todoCreatedMsg:
		cmds = append(cmds, m.loadTodos())
		cmds = append(cmds, ShowDefaultToast("Todo created", SuccessToast))

	case todoUpdatedMsg:
		cmds = append(cmds, m.loadTodos())
		cmds = append(cmds, ShowDefaultToast("Todo updated", SuccessToast))

	case todoDeletedMsg:
		cmds = append(cmds, m.loadTodos())
		cmds = append(cmds, ShowDefaultToast("Todo deleted", SuccessToast))

	case todoStatusChangedMsg:
		cmds = append(cmds, m.loadTodos())
		cmds = append(cmds, ShowDefaultToast(
			fmt.Sprintf("Todo status changed to %s", msg.newStatus),
			SuccessToast))

	case todoErrorMsg:
		cmds = append(cmds, ShowDefaultToast(msg.Error(), ErrorToast))

	case modalCloseMsg:
		m.showingModal = false
		if msg.reload {
			cmds = append(cmds, m.loadTodos())
		}
	}

	// Update list
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	// Update text input
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *TodoModel) View() string {
	if m.quitting {
		return ""
	}
	if m.showingModal && m.modalComponent != nil {
		return m.modalComponent.View()
	}

	// Header
	header := m.HeaderView()

	// Main list
	listView := lipgloss.NewStyle().Padding(styling.PaddingX).Render(m.list.View())

	// Input
	inputView := fmt.Sprintf(
		"\n%s\n",
		m.textInput.View(),
	)

	// Combine all views
	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		listView,
		inputView,
	)
}

func (m *TodoModel) HeaderView() string {
	filterButtons := []string{}
	smallLine := strings.Repeat("─", 2)
	for i := models.Open; i <= models.Archived; i++ {
		baseStyle := lipgloss.NewStyle()
		var textStyle lipgloss.Style
		if m.currentFilter == i {
			textStyle = styling.FocusedStyle
		} else {
			textStyle = lipgloss.NewStyle()
		}
		leftSide := fmt.Sprintf("%s[ ", smallLine)
		rightSide := fmt.Sprintf(" ]%s", smallLine)
		text := fmt.Sprintf("%d %s", i+1, i.String())
		filterButtons = append(filterButtons, baseStyle.Render(leftSide), textStyle.Render(text), baseStyle.Render(rightSide))
	}

	buttons := ""

	for _, b := range filterButtons {
		buttons += b
	}

	line := strings.Repeat("─", m.width-lipgloss.Width(buttons))

	return lipgloss.JoinHorizontal(lipgloss.Center, buttons, line)
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

// Message types
type todosLoadedMsg struct {
	todos []*models.Todo
}

type todoCreatedMsg struct {
	todo *models.Todo
}

type todoUpdatedMsg struct{}

type todoDeletedMsg struct{}

type todoStatusChangedMsg struct {
	newStatus models.Status
}

type todoErrorMsg struct {
	err error
}

func (e todoErrorMsg) Error() string {
	return e.err.Error()
}

type modalCloseMsg struct {
	reload bool
}

// Commands
func (m *TodoModel) loadTodos() tea.Cmd {
	return func() tea.Msg {
		var todos []*models.Todo
		var err error

		switch m.currentFilter {
		case models.Open:
			todos, err = m.service.GetOpenTodos()
		case models.Doing:
			todos, err = m.service.GetActiveTodos()
		case models.Done:
			todos, err = m.service.GetCompletedTodos()
		case models.Archived:
			todos, err = m.service.GetArchivedTodos()
		}

		if err != nil {
			return todoErrorMsg{err: err}
		}

		// Sort todos by priority (high to low)
		sort.Slice(todos, func(i, j int) bool {
			return todos[i].Priority > todos[j].Priority
		})

		return todosLoadedMsg{todos: todos}
	}
}

func (m *TodoModel) createTodoCmd(title, description string) tea.Cmd {
	return func() tea.Msg {
		todo, err := m.service.CreateTodo(title, description, models.Medium) // Default priority
		if err != nil {
			return todoErrorMsg{err: err}
		}
		return todoCreatedMsg{todo: todo}
	}
}

func (m *TodoModel) advanceTodoStatusCmd(todoID int64) tea.Cmd {
	return func() tea.Msg {
		todo, err := m.service.GetTodo(todoID)
		if err != nil {
			return todoErrorMsg{err: err}
		}

		var newStatus models.Status
		switch todo.Status {
		case models.Open:
			newStatus = models.Doing
			err = m.service.MarkAsDoing(todoID)
		case models.Doing:
			newStatus = models.Done
			err = m.service.MarkAsDone(todoID)
		case models.Done:
			newStatus = models.Archived
			err = m.service.ArchiveTodo(todoID)
		case models.Archived:
			newStatus = models.Open
			err = m.service.MarkAsOpen(todoID)
		}

		if err != nil {
			return todoErrorMsg{err: err}
		}

		return todoStatusChangedMsg{newStatus: newStatus}
	}
}

func (m *TodoModel) deleteTodoCmd(todoID int64) tea.Cmd {
	return func() tea.Msg {
		err := m.service.DeleteTodo(todoID)
		if err != nil {
			return todoErrorMsg{err: err}
		}
		return todoDeletedMsg{}
	}
}

func (m *TodoModel) showEditModalCmd(todo *models.Todo) tea.Cmd {
	return func() tea.Msg {
		m.showingModal = true
		m.modalComponent = NewTodoEditModal(todo, m.width, m.height, m.service)
		return nil
	}
}

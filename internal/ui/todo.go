package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type TodoItem struct {
	todo *models.Todo
}

func (i TodoItem) Title() string {
	return i.todo.Title
}

func (i TodoItem) Description() string {
	return i.todo.Description
}

func (i TodoItem) FilterValue() string {
	return i.todo.Title + " " + i.todo.Description
}

type TodoItemDelegate struct{}

func (d TodoItemDelegate) Height() int                             { return 1 }
func (d TodoItemDelegate) Spacing() int                            { return 0 }
func (d TodoItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d TodoItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(TodoItem)
	if !ok {
		return
	}
	width := m.Width() - 4

	// Left-aligned elements
	selected := styling.GetSelectedBlock(index == m.Index())
	priorityMarker := styling.GetStyledPriority(i.todo.Priority, true, false)
	title := styling.TextStyle.MarginRight(1).Width(50).Render(truncateString(i.Title(), 50))

	leftElementsWidth := lipgloss.Width(selected) + lipgloss.Width(priorityMarker) + lipgloss.Width(title)

	if leftElementsWidth >= width {
		widthAvailableForTitle := width - lipgloss.Width(selected) - 1
		shortTitle := styling.TextStyle.MarginRight(1).Width(widthAvailableForTitle).Render(truncateString(i.Title(), widthAvailableForTitle))
		row := lipgloss.JoinHorizontal(lipgloss.Center, selected, shortTitle)
		fmt.Fprint(w, row)
		return
	}

	// Right-aligned elements
	var rightElements []string

	// Add tags
	tags := ""
	for _, tag := range i.todo.Tags {
		tags += styling.GetStyledTag(tag)
	}
	if tags != "" {
		rightElements = append(rightElements, tags)
	}

	// Add due date if present
	dueDate := ""
	if i.todo.DueDate != nil {
		dueDate = styling.GetStyledDueDate(*i.todo.DueDate, i.todo.Priority)
		rightElements = append(rightElements, dueDate)
	}

	// Add updated at timestamp
	updatedAt := styling.GetStyledUpdatedAt(i.todo.UpdatedAt)
	rightElements = append(rightElements, updatedAt)

	// Join right elements
	rightContent := lipgloss.JoinHorizontal(lipgloss.Right, rightElements...)
	rightWidth := lipgloss.Width(rightContent)

	// Calculate space for description
	descriptionMaxWidth := width - leftElementsWidth - rightWidth - 2 // 2 for some padding

	// Truncate description if needed
	description := i.todo.Description
	if descriptionMaxWidth > 20 {
		description = truncateString(description, descriptionMaxWidth)
	} else {
		description = ""
	}

	styledDescription := styling.SubtextStyle.Width(descriptionMaxWidth).Render(description)

	// Assemble the row with left content taking remaining space and right content aligned to the right
	leftContent := lipgloss.JoinHorizontal(lipgloss.Left, selected, priorityMarker, title, styledDescription)

	// Join everything, ensuring right alignment for the right content
	row := lipgloss.NewStyle().Width(width).Render(
		lipgloss.JoinHorizontal(lipgloss.Left,
			leftContent,
			lipgloss.NewStyle().Width(width-lipgloss.Width(leftContent)).Align(lipgloss.Right).Render(rightContent),
		),
	)

	fmt.Fprint(w, row)
}

type TodoModel struct {
	service         *service.AppService
	tuiService      *service.TuiService
	list            list.Model
	width           int
	height          int
	quitting        bool
	modalComponent  tea.Model
	footer          tea.Model
	filterStatusBar tea.Model
}

func NewTodoModel(appService *service.AppService) *TodoModel {
	tuiService := service.NewTuiService()

	// Setup list
	todoList := list.New([]list.Item{}, TodoItemDelegate{}, 0, 0)
	todoList.Title = ""
	todoList.DisableQuitKeybindings()
	todoList.SetShowTitle(false)
	todoList.SetShowHelp(true)
	todoList.SetShowStatusBar(false)
	todoList.SetFilteringEnabled(true)

	footer := NewFooterModel(appService, tuiService)
	filterStatusBar := NewFilterStatusBar(tuiService)

	// Create model
	m := &TodoModel{
		service:         appService,
		tuiService:      tuiService,
		list:            todoList,
		footer:          footer,
		filterStatusBar: filterStatusBar,
	}

	todos, err := appService.GetActiveTodos()
	if err == nil && len(todos) > 0 {
		items := make([]list.Item, len(todos))
		for i, todo := range todos {
			items[i] = TodoItem{todo: todo}
		}
		m.list.SetItems(items)
	} else if err != nil {
		log.Error("Error pre-loading todos", "error", err)
	}
	return m
}

func (m *TodoModel) Init() tea.Cmd {
	return nil
}

func (m *TodoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.tuiService.ShouldShowModal() {
			// Handle modal
			m.modalComponent, cmd = m.modalComponent.Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(
			msg,
			m.tuiService.KeyMap.Quit,
		):
			if m.tuiService.CurrentView == service.NewView {
				m.tuiService.SwitchToListView()
			} else if m.list.FilterState() != 0 {
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			} else {
				m.quitting = true
				return m, tea.Quit
			}

		case key.Matches(msg, m.tuiService.KeyMap.SwitchPane):
			if m.tuiService.CurrentView == service.ListView {
				switch key := msg.String(); key {
				case "1", "2", "3", "4":
					m.tuiService.SelectFilter(key)
					return m, m.loadTodosCmd()
				case "5":
					m.tuiService.SwitchToNewTodoView()
					return m, nil
				}
			}

		case key.Matches(msg, m.tuiService.KeyMap.AdvanceStatus):
			// Advance todo status
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(TodoItem)
				return m, m.advanceTodoStatusCmd(item.todo.ID)
			}

		case key.Matches(msg, m.tuiService.KeyMap.Edit):
			// Edit selected todo
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(TodoItem)
				return m, m.showEditModalCmd(item.todo)
			}
		case key.Matches(msg, m.tuiService.KeyMap.New):
			// Create new Todo
			todo := &models.Todo{}
			return m, m.showEditModalCmd(todo)

		case key.Matches(msg, m.tuiService.KeyMap.Delete):
			// Delete selected todo
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(TodoItem)
				return m, m.showConfirmDeleteCmd(item.todo.ID)
			}
		case key.Matches(msg, m.tuiService.KeyMap.Archive):
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(TodoItem)
				return m, m.toggleArchiveCmd(item.todo.ID, item.todo.Archived)
			}
		case key.Matches(msg, m.tuiService.KeyMap.ToggleArchived):
			if m.tuiService.FilterState.Mode == service.AllFilter {
				m.tuiService.FilterState.IncludeArchived = !m.tuiService.FilterState.IncludeArchived
				return m, tea.Batch(
					m.loadTodosCmd(),
					ShowDefaultToast(
						fmt.Sprintf("Archived todos %s",
							map[bool]string{true: "shown", false: "hidden"}[m.tuiService.FilterState.IncludeArchived]),
						InfoToast),
				)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 3 // Title + top border
		footerHeight := 3 // Input + bottom padding

		m.list.SetSize(msg.Width, msg.Height-headerHeight-footerHeight)

		// Pass size to modal if active
		if m.tuiService.ShouldShowModal() && m.modalComponent != nil {
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
		m.tuiService.SwitchToListView()
		cmds = append(cmds, m.loadTodosCmd())
		cmds = append(cmds, ShowDefaultToast("Todo created", SuccessToast))

	case todoUpdatedMsg:
		cmds = append(cmds, m.loadTodosCmd())
		cmds = append(cmds, ShowDefaultToast("Todo updated", SuccessToast))

	case todoDeletedMsg:
		m.tuiService.SwitchToListView()
		cmds = append(cmds, m.loadTodosCmd())
		cmds = append(cmds, ShowDefaultToast("Todo deleted", SuccessToast))

	case todoStatusChangedMsg:
		cmds = append(cmds, m.loadTodosCmd())
		cmds = append(cmds, ShowDefaultToast(
			fmt.Sprintf("Todo status changed to %s", msg.newStatus),
			SuccessToast))

	case todoErrorMsg:
		cmds = append(cmds, ShowDefaultToast(msg.Error(), ErrorToast))

	case modalCloseMsg:
		m.tuiService.SwitchToListView()
		if msg.reload {
			cmds = append(cmds, m.loadTodosCmd())
		}

	case CreateTodoMsg:
		m.tuiService.SwitchToListView()
		cmds = append(cmds, m.createTodoCmd(msg.Title, msg.Priority))

	case showModalMsg:
		return m, tea.WindowSize()
	}

	// Update list
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	// Update text input
	if m.tuiService.CurrentView == service.NewView {
		m.footer, cmd = m.footer.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.tuiService.ShouldShowModal() && m.modalComponent != nil {
		m.modalComponent, cmd = m.modalComponent.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *TodoModel) View() string {
	if m.quitting {
		return ""
	}
	if m.tuiService.ShouldShowModal() && m.modalComponent != nil {
		return m.modalComponent.View()
	}

	header := m.HeaderView()
	filterBar := m.filterStatusBar.View()
	footer := m.footer.View()

	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	filterBarHeight := lipgloss.Height(filterBar)

	m.list.SetHeight(m.height - headerHeight - footerHeight - filterBarHeight - 4)

	// Main list
	listView := lipgloss.NewStyle().Width(m.width - 2).Padding(styling.Padding).Render(m.list.View())

	padding := lipgloss.NewStyle().Padding(1, 1)
	// Combine all views
	return padding.Render(lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		filterBar,
		listView,
		footer,
	))
}

func (m *TodoModel) HeaderView() string {
	var leftTabs []string

	for status := models.Open; status <= models.Done; status++ {
		isSelected := m.tuiService.FilterState.Mode == service.StatusFilter &&
			m.tuiService.FilterState.Status == status
		tab := styling.GetStyledStatus(status, isSelected, false)
		leftTabs = append(leftTabs, tab)
	}

	leftContent := lipgloss.JoinHorizontal(lipgloss.Center, leftTabs...)

	isAllSelected := m.tuiService.FilterState.Mode == service.AllFilter
	allTab := styling.GetStyledTagWithIndicator(4, "All", styling.Rosewater, isAllSelected, false)

	const minGap = 2
	availableWidth := m.width - 2 // -2 for padding
	leftWidth := lipgloss.Width(leftContent)
	rightWidth := lipgloss.Width(allTab)

	if leftWidth+minGap+rightWidth >= availableWidth {
		return lipgloss.JoinHorizontal(lipgloss.Center, leftContent, allTab)
	}

	spacerWidth := availableWidth - leftWidth - rightWidth
	spacer := strings.Repeat(" ", spacerWidth)

	return lipgloss.JoinHorizontal(lipgloss.Center, leftContent, spacer, allTab)
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

type todoCreatedMsg struct{}

type todoUpdatedMsg struct{}

type showModalMsg struct{}

type todoStatusChangedMsg struct {
	newStatus string
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
func (m *TodoModel) loadTodosCmd() tea.Cmd {
	return func() tea.Msg {
		todos, err := m.service.GetFilteredTodos(
			m.tuiService.FilterState.Mode,
			m.tuiService.FilterState.Status,
			m.tuiService.FilterState.Tag,
			m.tuiService.FilterState.IncludeArchived,
		)

		if err != nil {
			return todoErrorMsg{err: err}
		}

		return todosLoadedMsg{todos: todos}
	}
}

func (m *TodoModel) createTodoCmd(title string, priority models.Priority) tea.Cmd {
	return func() tea.Msg {
		err := m.service.CreateTodo(title, "", priority, []string{})
		if err != nil {
			return todoErrorMsg{err: err}
		}
		return todoCreatedMsg{}
	}
}

func (m *TodoModel) advanceTodoStatusCmd(todoID int64) tea.Cmd {
	return func() tea.Msg {
		newStatus, err := m.service.AdvanceStatus(todoID)
		if err != nil {
			return todoErrorMsg{err: err}
		}

		return todoStatusChangedMsg{newStatus: newStatus.String()}
	}
}

func (m *TodoModel) showEditModalCmd(todo *models.Todo) tea.Cmd {
	return func() tea.Msg {
		m.tuiService.SwitchToEditTodoView()
		m.modalComponent = NewTodoEditModal(todo, m.width, m.height, m.service, m.tuiService)
		return showModalMsg{}
	}
}

func (m *TodoModel) showConfirmDeleteCmd(todoID int64) tea.Cmd {
	return func() tea.Msg {
		m.tuiService.SwitchToConfirmDeleteView()
		m.modalComponent = NewConfirmDeleteModal(m.service, m.tuiService, todoID)
		return showModalMsg{}
	}
}

func CloseModalCmd(reload bool) tea.Cmd {
	return func() tea.Msg {
		return modalCloseMsg{reload: reload}
	}
}

func (m *TodoModel) toggleArchiveCmd(todoID int64, isArchived bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		if isArchived {
			err = m.service.UnarchiveTodo(todoID)
		} else {
			err = m.service.ArchiveTodo(todoID)
		}

		if err != nil {
			return todoErrorMsg{err: err}
		}

		action := "archived"
		if isArchived {
			action = "unarchived"
		}

		return todoStatusChangedMsg{
			newStatus: fmt.Sprintf("Todo %s", action),
		}
	}
}

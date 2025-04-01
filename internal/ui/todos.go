package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type TodosModel struct {
	service        *service.AppService
	tuiService     *service.TuiService
	translator     *i18n.TranslationService
	list           list.Model
	width          int
	height         int
	quitting       bool
	modalComponent tea.Model
	footer         tea.Model
}

func NewTodoModel(appService *service.AppService, translationService *i18n.TranslationService) *TodosModel {
	tuiService := service.NewTuiService()

	// Setup list
	todoList := list.New([]list.Item{}, TodoModel{translator: translationService}, 0, 0)
	todoList.Title = ""
	todoList.DisableQuitKeybindings()
	todoList.SetShowTitle(false)
	todoList.SetShowHelp(false)
	todoList.SetShowStatusBar(false)
	todoList.SetFilteringEnabled(true)

	footer := NewFooterModel(appService, tuiService, translationService)

	// Create model
	m := &TodosModel{
		service:    appService,
		tuiService: tuiService,
		translator: translationService,
		list:       todoList,
		footer:     footer,
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

func (m *TodosModel) Init() tea.Cmd {
	return nil
}

func (m *TodosModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.tuiService.CurrentView == service.AddEditView {
				m.tuiService.SwitchToListView()
			} else if m.list.FilterState() != 0 {
				m.tuiService.RemoveNameFilter()
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			} else {
				m.quitting = true
				return m, tea.Quit
			}

		case key.Matches(msg, m.tuiService.KeyMap.Filter):
			m.tuiService.ActivateNameFilter()
			m.list, cmd = m.list.Update(msg)
			return m, cmd

		case key.Matches(msg, m.tuiService.KeyMap.SwitchPane):
			if m.tuiService.CurrentView == service.ListView {
				switch key := msg.String(); key {
				case "1", "2", "3", "4":
					m.tuiService.SwitchPane(key)
					return m, m.loadTodosCmd()
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
				key := fmt.Sprintf("toast.filter_archived_%s",
					map[bool]string{true: "shown", false: "hidden"}[m.tuiService.FilterState.IncludeArchived])
				return m, tea.Batch(
					m.loadTodosCmd(),
					ShowDefaultToast(
						m.translator.T(key),
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

		m.footer, cmd = m.footer.Update(msg)
		cmds = append(cmds, cmd)

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
		cmds = append(cmds, ShowDefaultToast(m.translator.T("toast.todo_created"), SuccessToast))

	case todoUpdatedMsg:
		cmds = append(cmds, m.loadTodosCmd())
		cmds = append(cmds, ShowDefaultToast(m.translator.T("toast.todo_updated"), SuccessToast))

	case todoDeletedMsg:
		m.tuiService.SwitchToListView()
		cmds = append(cmds, m.loadTodosCmd())
		cmds = append(cmds, ShowDefaultToast(m.translator.T("toast.todo_deleted"), SuccessToast))

	case todoStatusChangedMsg:
		cmds = append(cmds, m.loadTodosCmd())
		cmds = append(cmds, ShowDefaultToast(
			m.translator.Tf("toast.status_changed", map[string]interface{}{"Status": m.translator.T(msg.newStatus)}),
			SuccessToast))

	case todoToggleArchived:
		cmds = append(cmds, m.loadTodosCmd())
		cmds = append(cmds, ShowDefaultToast(
			m.translator.T(fmt.Sprintf("toast.%s", msg.action)),
			SuccessToast))

	case todoErrorMsg:
		cmds = append(cmds, ShowDefaultToast(msg.Error(), ErrorToast))

	case modalCloseMsg:
		m.tuiService.SwitchToListView()
		if msg.reload {
			cmds = append(cmds, m.loadTodosCmd())
		}

	case showModalMsg:
		return m, tea.WindowSize()
	}

	// Update list
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	// Update text input
	if m.tuiService.CurrentView == service.ListView {
		m.footer, cmd = m.footer.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.tuiService.ShouldShowModal() && m.modalComponent != nil {
		m.modalComponent, cmd = m.modalComponent.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *TodosModel) View() string {
	if m.quitting {
		return ""
	}
	if m.tuiService.ShouldShowModal() && m.modalComponent != nil {
		return m.modalComponent.View()
	}

	header := m.HeaderView()
	footer := m.footer.View()

	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)

	m.list.SetHeight(m.height - headerHeight - footerHeight - 4)

	// Main list
	listView := lipgloss.NewStyle().Width(m.width - 2).Padding(styling.Padding).Render(m.list.View())

	padding := lipgloss.NewStyle().Padding(1, 1)
	// Combine all views
	return padding.Render(lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		listView,
		footer,
	))
}

func (m *TodosModel) HeaderView() string {
	var leftTabs []string

	for status := models.Open; status <= models.Done; status++ {
		isSelected := m.tuiService.FilterState.Mode == service.StatusFilter &&
			m.tuiService.FilterState.Status == status
		translatedStatus := m.translator.T(status.String())
		tab := styling.GetStyledStatus(translatedStatus, status, isSelected, false)
		leftTabs = append(leftTabs, tab)
	}

	leftContent := lipgloss.JoinHorizontal(lipgloss.Center, leftTabs...)

	isAllSelected := m.tuiService.FilterState.Mode == service.AllFilter
	allTab := styling.GetStyledTagWithIndicator(4, m.translator.T("filter.all"), styling.Rosewater, isAllSelected, false)

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

// ===========================================================================
// Helpers
// ===========================================================================
func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

// ===========================================================================
// Message Types
// ===========================================================================
type todoToggleArchived struct {
	action string
}
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

// ===========================================================================
// Commands
// ===========================================================================
func (m *TodosModel) loadTodosCmd() tea.Cmd {
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

func (m *TodosModel) advanceTodoStatusCmd(todoID int64) tea.Cmd {
	return func() tea.Msg {
		newStatus, err := m.service.AdvanceStatus(todoID)
		if err != nil {
			return todoErrorMsg{err: err}
		}

		return todoStatusChangedMsg{newStatus: newStatus.String()}
	}
}

func (m *TodosModel) showEditModalCmd(todo *models.Todo) tea.Cmd {
	return func() tea.Msg {
		m.tuiService.SwitchToEditTodoView()
		m.modalComponent = NewTodoEditModal(todo, m.width, m.height, m.service, m.tuiService, m.translator)
		return showModalMsg{}
	}
}

func (m *TodosModel) showConfirmDeleteCmd(todoID int64) tea.Cmd {
	return func() tea.Msg {
		m.tuiService.SwitchToConfirmDeleteView()
		m.modalComponent = NewConfirmDeleteModal(m.service, m.tuiService, m.translator, todoID)
		return showModalMsg{}
	}
}

func CloseModalCmd(reload bool) tea.Cmd {
	return func() tea.Msg {
		return modalCloseMsg{reload: reload}
	}
}

func (m *TodosModel) toggleArchiveCmd(todoID int64, isArchived bool) tea.Cmd {
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

		return todoToggleArchived{
			action: action,
		}
	}
}

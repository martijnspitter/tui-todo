package ui

import (
	"fmt"

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
	header         tea.Model
	today          tea.Model
}

func NewTodosModel(appService *service.AppService, translationService *i18n.TranslationService) *TodosModel {
	tuiService := service.NewTuiService()

	// Setup list
	todoList := list.New([]list.Item{}, TodoModel{translator: translationService, tuiService: tuiService}, 0, 0)
	todoList.Title = ""
	todoList.DisableQuitKeybindings()
	todoList.SetShowTitle(false)
	todoList.SetShowHelp(false)
	todoList.SetShowStatusBar(false)
	todoList.SetFilteringEnabled(true)

	footer := NewFooterModel(appService, tuiService, translationService)
	header := NewHeaderModel(tuiService, translationService)
	today := NewTodayModel(appService, tuiService, translationService)

	// Create model
	m := &TodosModel{
		service:    appService,
		tuiService: tuiService,
		translator: translationService,
		list:       todoList,
		footer:     footer,
		header:     header,
		today:      today,
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
	case LoadTodosMsg:
		log.Debug("LOAD TODOS SYNC")
		return m, m.loadTodosCmd()
	case tea.KeyMsg:
		if m.tuiService.ShouldShowModal() {
			// Handle modal
			m.modalComponent, cmd = m.modalComponent.Update(msg)
			return m, cmd
		}

		switch {
		case !key.Matches(msg, m.tuiService.KeyMap.Quit) && m.tuiService.ShowConfirmQuit:
			m.tuiService.ToggleShowConfirmQuit()

		case key.Matches(
			msg,
			m.tuiService.KeyMap.Quit,
		):
			if m.tuiService.CurrentView == service.AddEditModal {
				m.tuiService.SwitchToListView()
			} else if m.list.FilterState() != 0 {
				m.tuiService.RemoveNameFilter()
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			} else if !m.tuiService.ShowConfirmQuit {
				m.tuiService.ToggleShowConfirmQuit()
			} else {
				m.quitting = true
				return m, tea.Quit
			}

		case key.Matches(msg, m.tuiService.KeyMap.TagFilter):
			if !m.tuiService.FilterState.IsFilterActive {
				m.tuiService.ActivateTagFilter()
				m.list.ResetFilter()
				filterKeyMsg := tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{'/'},
				}
				m.list, cmd = m.list.Update(filterKeyMsg)
				return m, cmd
			}

		case key.Matches(msg, m.tuiService.KeyMap.Filter):
			if !m.tuiService.FilterState.IsFilterActive {
				m.tuiService.ActivateTitleFilter()
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}

		case key.Matches(msg, m.tuiService.KeyMap.SwitchPane):
			if !m.tuiService.ShouldShowModal() {
				m.tuiService.SwitchPane(msg.String())
				return m, m.loadTodosCmd()
			}

		case key.Matches(msg, m.tuiService.KeyMap.AdvanceStatus):
			// Advance todo status
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(*TodoItem)
				return m, m.advanceTodoStatusCmd(item.todo.ID)
			}

		case key.Matches(msg, m.tuiService.KeyMap.Edit):
			// Edit selected todo
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(*TodoItem)
				return m, m.showEditModalCmd(item.todo)
			}
		case key.Matches(msg, m.tuiService.KeyMap.New):
			// Create new Todo
			todo := &models.Todo{}
			return m, m.showEditModalCmd(todo)

		case key.Matches(msg, m.tuiService.KeyMap.Delete):
			// Delete selected todo
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(*TodoItem)
				return m, m.showConfirmDeleteCmd(item.todo.ID)
			}
		case key.Matches(msg, m.tuiService.KeyMap.Archive):
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(*TodoItem)
				return m, m.toggleArchiveCmd(item.todo.ID, item.todo.Archived)
			}
		case key.Matches(msg, m.tuiService.KeyMap.ToggleArchived):
			if m.tuiService.CurrentView == service.AllPane {
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

		case key.Matches(msg, m.tuiService.KeyMap.About):
			return m, m.showAboutModalCmd()
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
			items[i] = &TodoItem{todo: todo, tuiService: m.tuiService}
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

	case TodoErrorMsg:
		cmds = append(cmds, ShowDefaultToast(m.translator.T(msg.Error()), ErrorToast))

	case modalCloseMsg:
		m.tuiService.SwitchToListView()
		if msg.reload {
			cmds = append(cmds, m.loadTodosCmd())
		}

	case UpdateCheckCompletedMsg:
		if msg.ForceUpdate {
			return m, m.showUpdateModalCmd()
		}

	case showModalMsg:
		return m, tea.WindowSize()
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	m.footer, cmd = m.footer.Update(msg)
	cmds = append(cmds, cmd)

	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	m.today, cmd = m.today.Update(msg)
	cmds = append(cmds, cmd)

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

	header := m.header.View()
	footer := m.footer.View()

	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)

	m.list.SetHeight(m.height - headerHeight - footerHeight - 3)

	// Main list
	listView := lipgloss.NewStyle().Width(m.width - 2).Padding(styling.Padding).Render(m.list.View())
	if m.tuiService.CurrentView == service.TodayPane {
		listView = m.today.View()
	} else if len(m.list.Items()) == 0 {
		listView = EmptyStateView(m.translator, m.width, m.height-headerHeight-footerHeight-1)
	}

	padding := lipgloss.NewStyle().Padding(1, 1)
	// Combine all views
	return padding.Render(lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		listView,
		footer,
	))
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

type TodoErrorMsg struct {
	err error
}

func (e TodoErrorMsg) Error() string {
	return e.err.Error()
}

type modalCloseMsg struct {
	reload bool
}

type UpdateCheckCompletedMsg struct {
	ForceUpdate bool
}

type LoadTodosMsg struct{}

// ===========================================================================
// Commands
// ===========================================================================
func (m *TodosModel) loadTodosCmd() tea.Cmd {
	return func() tea.Msg {
		if m.tuiService.CurrentView == service.TodayPane {
			return GetTodayDataMsg{}
		}

		todos, err := m.service.GetFilteredTodos(
			m.tuiService.CurrentView,
			m.tuiService.FilterState.IncludeArchived,
		)

		if err != nil {
			return TodoErrorMsg{err: err}
		}

		return todosLoadedMsg{todos: todos}
	}
}

func (m *TodosModel) advanceTodoStatusCmd(todoID int64) tea.Cmd {
	return func() tea.Msg {
		newStatus, err := m.service.AdvanceStatus(todoID)
		if err != nil {
			return TodoErrorMsg{err: err}
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
			return TodoErrorMsg{err: err}
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

func (m *TodosModel) showUpdateModalCmd() tea.Cmd {
	return func() tea.Msg {
		m.tuiService.SwitchToUpdateModalView()
		m.modalComponent = NewUpdateModal(
			m.width,
			m.height,
			m.service,
			m.tuiService,
			m.translator,
		)
		return showModalMsg{}
	}
}

func (m *TodosModel) showAboutModalCmd() tea.Cmd {
	return func() tea.Msg {
		m.tuiService.SwitchToAboutModalView()
		m.modalComponent = NewAboutModal(
			m.width,
			m.height,
			m.service,
			m.tuiService,
			m.translator,
		)
		return showModalMsg{}
	}
}

func InitCmd() tea.Msg {
	return LoadTodosMsg{}
}

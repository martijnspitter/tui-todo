package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type TodosModel struct {
	service    *service.AppService
	tuiService *service.TuiService
	translator *i18n.TranslationService
	list       list.Model
	width      int
	height     int
}

func NewTodosModel(service *service.AppService, tuiService *service.TuiService, translator *i18n.TranslationService) *TodosModel {
	// Setup list
	todoList := list.New([]list.Item{}, TodoModel{translator: translator, tuiService: tuiService}, 0, 0)
	todoList.Title = ""
	todoList.DisableQuitKeybindings()
	todoList.SetShowTitle(false)
	todoList.SetShowHelp(false)
	todoList.SetShowStatusBar(false)
	todoList.SetFilteringEnabled(true)

	return &TodosModel{
		service:    service,
		tuiService: tuiService,
		translator: translator,
		list:       todoList,
	}
}

func (m *TodosModel) Init() tea.Cmd {
	return nil
}

func (m *TodosModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
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
		case key.Matches(msg, m.tuiService.KeyMap.AdvanceStatus):
			// Advance todo status
			if m.shouldAllowTodoCrud() {
				item := m.list.SelectedItem().(*TodoItem)
				return m, m.advanceTodoStatusCmd(item.todo.ID)
			}
		case key.Matches(msg, m.tuiService.KeyMap.Edit):
			// Edit selected todo
			if m.shouldAllowTodoCrud() {
				item := m.list.SelectedItem().(*TodoItem)
				return m, m.showEditModalCmd(item.todo)
			}
		case key.Matches(msg, m.tuiService.KeyMap.Delete):
			// Delete selected todo
			if m.shouldAllowTodoCrud() {
				item := m.list.SelectedItem().(*TodoItem)
				return m, m.showConfirmDeleteCmd(item.todo.ID)
			}
		case key.Matches(msg, m.tuiService.KeyMap.Archive):
			if m.shouldAllowTodoCrud() {
				item := m.list.SelectedItem().(*TodoItem)
				return m, m.toggleArchiveCmd(item.todo.ID, item.todo.Archived)
			}
		case key.Matches(msg, m.tuiService.KeyMap.BlockTodo):
			// Block/unblock todo
			if m.shouldAllowTodoCrud() {
				item := m.list.SelectedItem().(*TodoItem)
				isCurrentlyBlocked := item.todo.Status == models.Blocked
				return m, m.blockTodoCmd(item.todo.ID, isCurrentlyBlocked)
			}
		case key.Matches(msg, m.tuiService.KeyMap.New):
			// Create new Todo
			todo := &models.Todo{}
			return m, m.showEditModalCmd(todo)
		}
	case RemoveFilterMsg:
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		headerHeight := 3 // Title + top border
		footerHeight := 3 // Input + bottom padding
		m.width = msg.Width
		m.height = msg.Height - headerHeight - footerHeight

		m.list.SetSize(msg.Width, m.height)
	case todosLoadedMsg:
		items := make([]list.Item, len(msg.todos))
		for i, todo := range msg.todos {
			items[i] = &TodoItem{todo: todo, tuiService: m.tuiService}
		}
		cmd := m.list.SetItems(items)
		cmds = append(cmds, cmd)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *TodosModel) View() string {
	listView := lipgloss.NewStyle().Width(m.width - 2).Padding(styling.Padding).Render(m.list.View())
	if len(m.list.Items()) == 0 {
		if m.tuiService.CurrentView == service.AllPane || m.tuiService.CurrentView == service.BlockedPane {
			listView = EmptyNothingFoundView(m.translator, m.width, m.height)
		} else {
			listView = EmptySuccessStateView(m.translator, m.width, m.height)
		}
	}

	return listView
}

// ===========================================================================
// Helpers
// ===========================================================================
func (m *TodosModel) shouldAllowTodoCrud() bool {
	return m.list.SelectedItem() != nil && m.tuiService.CurrentView != service.TodayPane && m.tuiService.CurrentView != service.TagsPane
}

func (m *TodosModel) SetHeight(height int) {
	m.height = height
	m.list.SetHeight(height)
}

// ===========================================================================
// Commands
// ===========================================================================
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
		modalComponent := NewTodoEditModal(todo, m.width, m.height, m.service, m.tuiService, m.translator)
		return showModalMsg{
			modal: modalComponent,
		}
	}
}

func (m *TodosModel) showConfirmDeleteCmd(todoID int64) tea.Cmd {
	return func() tea.Msg {
		m.tuiService.SwitchToConfirmDeleteView()
		modalComponent := NewConfirmDeleteModal(m.service, m.tuiService, m.translator, todoID, false)
		return showModalMsg{
			modal: modalComponent,
		}
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

func (m *TodosModel) blockTodoCmd(todoID int64, isCurrentlyBlocked bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		var newStatus models.Status

		if isCurrentlyBlocked {
			// If currently blocked, unblock by setting to Open
			err = m.service.MarkAsOpen(todoID)
			newStatus = models.Open
		} else {
			// If not blocked, block it
			err = m.service.MarkAsBlocked(todoID)
			newStatus = models.Blocked
		}

		if err != nil {
			return TodoErrorMsg{err: err}
		}

		return todoStatusChangedMsg{
			newStatus: newStatus.String(),
		}
	}
}

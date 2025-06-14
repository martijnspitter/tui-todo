package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
)

type MainModel struct {
	service        *service.AppService
	tuiService     *service.TuiService
	translator     *i18n.TranslationService
	width          int
	height         int
	quitting       bool
	modalComponent tea.Model
	footer         tea.Model
	header         tea.Model
	today          tea.Model
	todos          tea.Model
	tags           tea.Model
}

func NewMainModel(appService *service.AppService, translationService *i18n.TranslationService) *MainModel {
	tuiService := service.NewTuiService()

	footer := NewFooterModel(appService, tuiService, translationService)
	header := NewHeaderModel(tuiService, translationService)
	today := NewTodayModel(appService, tuiService, translationService)
	todos := NewTodosModel(appService, tuiService, translationService)
	tags := NewTagsModel(appService, tuiService, translationService)

	// Create model
	m := &MainModel{
		service:    appService,
		tuiService: tuiService,
		translator: translationService,
		todos:      todos,
		footer:     footer,
		header:     header,
		today:      today,
		tags:       tags,
	}

	return m
}

func (m *MainModel) Init() tea.Cmd {
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case LoadTodosMsg:
		return m, m.loadTodosCmd()
	case LoadTagsMsg:
		return m, m.loadTagsCmd()
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
			if m.tuiService.CurrentView == service.AddEditTodoModal || m.tuiService.CurrentView == service.AddEditTagModal {
				m.tuiService.SwitchToListView()
			} else if m.tuiService.FilterState.IsFilterActive {
				m.tuiService.RemoveNameFilter()
				cmd := RemoveFilterCmd()
				cmds = append(cmds, cmd)
			} else if !m.tuiService.ShowConfirmQuit {
				m.tuiService.ToggleShowConfirmQuit()
			} else {
				m.quitting = true
				return m, tea.Quit
			}

		case key.Matches(msg, m.tuiService.KeyMap.SwitchPane):
			if !m.tuiService.ShouldShowModal() {
				m.tuiService.SwitchPane(msg.String())
				return m, m.loadTodosCmd()
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

		m.footer, cmd = m.footer.Update(msg)
		cmds = append(cmds, cmd)

		// Pass size to modal if active
		if m.tuiService.ShouldShowModal() && m.modalComponent != nil {
			var cmd tea.Cmd
			m.modalComponent, cmd = m.modalComponent.Update(msg)
			cmds = append(cmds, cmd)
		}

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
	case tagModalCloseMsg:
		m.tuiService.SwitchToTagsView()
		if msg.reload {
			cmds = append(cmds, m.loadTagsCmd())
		}

	case UpdateCheckCompletedMsg:
		if msg.ForceUpdate {
			return m, m.showUpdateModalCmd()
		}

	case showModalMsg:
		m.modalComponent = msg.modal
		return m, tea.WindowSize()
	}

	m.footer, cmd = m.footer.Update(msg)
	cmds = append(cmds, cmd)

	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	m.today, cmd = m.today.Update(msg)
	cmds = append(cmds, cmd)

	m.todos, cmd = m.todos.Update(msg)
	cmds = append(cmds, cmd)

	m.tags, cmd = m.tags.Update(msg)
	cmds = append(cmds, cmd)

	if m.tuiService.ShouldShowModal() && m.modalComponent != nil {
		m.modalComponent, cmd = m.modalComponent.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *MainModel) View() string {
	if m.quitting {
		return ""
	}
	if m.tuiService.ShouldShowModal() && m.modalComponent != nil {
		return m.modalComponent.View()
	}

	header := m.header.View()
	footer := m.footer.View()
	todos := m.todos.View()
	tags := m.tags.View()

	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)

	contentHeight := m.height - headerHeight - footerHeight - 3
	if todosModel, ok := m.todos.(*TodosModel); ok {
		todosModel.SetHeight(contentHeight)
	}
	if tagsModel, ok := m.tags.(*TagsModel); ok {
		tagsModel.SetHeight(contentHeight)
	}

	// Main list
	listView := ""
	if m.tuiService.CurrentView == service.TodayPane {
		listView = m.today.View()
	} else if m.tuiService.CurrentView == service.TagsPane {
		listView = tags
	} else {
		listView = todos
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
	if length <= 0 {
		return ""
	}
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
type tagsLoadedMsg struct {
	tags []*models.Tag
}

type todoCreatedMsg struct{}

type todoUpdatedMsg struct{}

type showModalMsg struct {
	modal tea.Model
}

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
type tagModalCloseMsg struct {
	reload bool
}

type UpdateCheckCompletedMsg struct {
	ForceUpdate bool
}

type LoadTodosMsg struct{}
type LoadTagsMsg struct{}

type RemoveFilterMsg struct{}

// ===========================================================================
// Commands
// ===========================================================================
func (m *MainModel) loadTodosCmd() tea.Cmd {
	return func() tea.Msg {
		if m.tuiService.CurrentView == service.TodayPane {
			return GetTodayDataMsg{}
		}

		if m.tuiService.CurrentView == service.TagsPane {
			return LoadTagsMsg{}
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

func (m *MainModel) loadTagsCmd() tea.Cmd {
	return func() tea.Msg {
		tags, err := m.service.GetAllTags()
		if err != nil {
			return TodoErrorMsg{err: err}
		}

		return tagsLoadedMsg{tags: tags}
	}
}

func CloseModalCmd(reload bool) tea.Cmd {
	return func() tea.Msg {
		return modalCloseMsg{reload: reload}
	}
}

func (m *MainModel) showUpdateModalCmd() tea.Cmd {
	return func() tea.Msg {
		m.tuiService.SwitchToUpdateModalView()
		modal := NewUpdateModal(
			m.width,
			m.height,
			m.service,
			m.tuiService,
			m.translator,
		)
		return showModalMsg{
			modal: modal,
		}
	}
}

func (m *MainModel) showAboutModalCmd() tea.Cmd {
	return func() tea.Msg {
		m.tuiService.SwitchToAboutModalView()
		modal := NewAboutModal(
			m.width,
			m.height,
			m.service,
			m.tuiService,
			m.translator,
		)
		return showModalMsg{
			modal: modal,
		}
	}
}

func InitTodosCmd() tea.Cmd {
	return func() tea.Msg {
		return LoadTodosMsg{}
	}
}

func InitTagsCmd() tea.Cmd {
	return func() tea.Msg {
		return LoadTagsMsg{}
	}
}

func RemoveFilterCmd() tea.Cmd {
	return func() tea.Msg {
		return RemoveFilterMsg{}
	}
}

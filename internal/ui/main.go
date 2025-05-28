package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

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
}

func NewMainModel(appService *service.AppService, translationService *i18n.TranslationService) *MainModel {
	tuiService := service.NewTuiService()

	footer := NewFooterModel(appService, tuiService, translationService)
	header := NewHeaderModel(tuiService, translationService)
	today := NewTodayModel(appService, tuiService, translationService)
	todos := NewTodosModel(appService, tuiService, translationService)

	// Create model
	m := &MainModel{
		service:    appService,
		tuiService: tuiService,
		translator: translationService,
		todos:      todos,
		footer:     footer,
		header:     header,
		today:      today,
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

	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)

	contentHeight := m.height - headerHeight - footerHeight - 3
	m.todos.(*TodosModel).SetHeight(contentHeight)

	// Main list
	listView := ""
	if m.tuiService.CurrentView == service.TodayPane {
		listView = m.today.View()
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
	log.Debug("s", s, "length", length)
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

type UpdateCheckCompletedMsg struct {
	ForceUpdate bool
}

type LoadTodosMsg struct{}

type RemoveFilterMsg struct{}

// ===========================================================================
// Commands
// ===========================================================================
func (m *MainModel) loadTodosCmd() tea.Cmd {
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

func CloseModalCmd(reload bool) tea.Cmd {
	return func() tea.Msg {
		return modalCloseMsg{reload: reload}
	}
}

func (m *MainModel) showUpdateModalCmd() tea.Cmd {
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

func (m *MainModel) showAboutModalCmd() tea.Cmd {
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

func RemoveFilterCmd() tea.Cmd {
	return func() tea.Msg {
		return RemoveFilterMsg{}
	}
}

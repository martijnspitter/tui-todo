package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
	"github.com/martijnspitter/tui-todo/internal/theme"
)

type TodayDashboardModel struct {
	service    *service.AppService
	tuiService *service.TuiService
	translator *i18n.TranslationService
	width      int
	height     int

	viewport viewport.Model
	ready    bool

	// Section data
	highPriorityTasks []*models.Todo
	dueTodayTasks     []*models.Todo
	inProgressTasks   []*models.Todo
	overdueTasks      []*models.Todo
	upcomingTasks     []*models.Todo

	// Stats
	completedTasksCount int
	totalTasksCount     int
	formattedTimeSpent  string

	// Currently selected section and item
	activeSection         int
	selectedItemInSection int64
}

func NewTodayModel(service *service.AppService, tuiService *service.TuiService, translator *i18n.TranslationService) *TodayDashboardModel {
	return &TodayDashboardModel{
		service:    service,
		tuiService: tuiService,
		translator: translator,
	}
}

func (m *TodayDashboardModel) Init() tea.Cmd {
	return nil
}

func (m *TodayDashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case GetTodayDataMsg:
		cmd = m.GetTodayDataCmd()
		return m, cmd
	case GetCompletionStats:
		cmd = m.GetCompletionStatsCmd()
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 5

		if !m.ready {
			// Initialize viewport when we first get a window size
			m.viewport = viewport.New(m.width, m.height)
			m.viewport.Style = lipgloss.NewStyle().Padding(0)
			m.ready = true
		} else {
			// Update viewport size
			m.viewport.Width = m.width
			m.viewport.Height = m.height + 1
		}

		// Update content if we have any
		contentStr := m.renderDashboard()
		m.viewport.SetContent(contentStr)

	case TodayDataUpdatedMsg:
		// When data is updated, update the viewport content
		if m.ready {
			contentStr := m.renderDashboard()
			m.viewport.SetContent(contentStr)
		}
	}

	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *TodayDashboardModel) View() string {
	if !m.ready {
		return m.loadingView()
	}

	// Show the viewport with scroll indicators if necessary
	scrollIndicator := ""
	if m.viewport.ScrollPercent() < 1.0 {
		scrollIndicator = styling.SubtextStyle.Render("\n↓ Scroll for more")
	}

	return m.viewport.View() + scrollIndicator
}

// ===========================================================================
// Helpers
// ===========================================================================
func (m *TodayDashboardModel) renderDashboard() string {
	// Calculate available width for content inside the modal
	modalMaxWidth := 150
	modalWidth := min(m.width-10, modalMaxWidth)
	contentWidth := modalWidth - 6 // Account for padding
	mainBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Mauve).
		Padding(1, 2).
		Width(modalWidth)

	// Progress bar and overview
	overviewTitle := styling.TextStyle.Bold(true).Render(m.translator.T("today_overview_title"))
	progressBar := m.renderProgressBar()

	timeSpentText := m.translator.Tf("ui.t_time_spent", map[string]interface{}{"Time": m.formattedTimeSpent})
	timeSpent := styling.GetTimeSpend(timeSpentText)

	overviewContent := lipgloss.JoinVertical(
		lipgloss.Center,
		overviewTitle,
		"",
		progressBar,
		"",
		timeSpent,
	)
	overviewBox := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		MarginBottom(1).
		Render(overviewContent)

	// If no data is loaded yet, show loading state
	if m.allTodosEmpty() {
		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Green).
			Padding(1, 2).
			Width(contentWidth).
			Render(lipgloss.JoinVertical(
				lipgloss.Left,
				EmptyStateView(m.translator, contentWidth, 13),
			))
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			overviewBox,
			emptyBox,
		)

		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			mainBox.Render(content),
		)
	}

	// High Priority Tasks Section
	highPrioTitle := styling.TextStyle.
		Bold(true).
		Foreground(theme.ErrorColor).
		Render(m.translator.Tf("today_high_prio",
			map[string]interface{}{"count": len(m.highPriorityTasks)}))

	highPrioTasks := m.renderTasks(m.highPriorityTasks, contentWidth)
	highPrioBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.ErrorColor).
		Padding(1, 2).
		Width(contentWidth).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			highPrioTitle,
			"",
			highPrioTasks,
		))

	// Overdue Tasks Section
	overdueTitle := styling.TextStyle.
		Bold(true).
		Foreground(theme.WarningColor).
		Render(m.translator.Tf("today_over_due",
			map[string]interface{}{"count": len(m.overdueTasks)}))

	overdueTasks := m.renderTasks(m.overdueTasks, contentWidth)
	overdueBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.WarningColor).
		Padding(1, 2).
		Width(contentWidth).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			overdueTitle,
			"",
			overdueTasks,
		))

	// Due Today Tasks Section
	dueTodayTitle := styling.TextStyle.
		Bold(true).
		Foreground(theme.Yellow).
		Render(m.translator.Tf("today_due_today",
			map[string]interface{}{"count": len(m.dueTodayTasks)}))

	dueTodayTasks := m.renderTasks(m.dueTodayTasks, contentWidth)
	dueTodayBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Yellow).
		Padding(1, 2).
		Width(contentWidth).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			dueTodayTitle,
			"",
			dueTodayTasks,
		))

	// In Progress Tasks Section
	inProgressTitle := styling.TextStyle.
		Bold(true).
		Foreground(theme.DoingStatusColor).
		Render(m.translator.Tf("today_in_progress",
			map[string]interface{}{"count": len(m.inProgressTasks)}))

	inProgressTasks := m.renderTasks(m.inProgressTasks, contentWidth)
	inProgressBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.DoingStatusColor).
		Padding(1, 2).
		Width(contentWidth).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			inProgressTitle,
			"",
			inProgressTasks,
		))

	// Coming Up Tasks Section
	upcomingTitle := styling.TextStyle.
		Bold(true).
		Foreground(theme.InfoColor).
		Render(m.translator.Tf("today_coming_up",
			map[string]interface{}{"count": len(m.upcomingTasks)}))

	upcomingTasks := m.renderTasks(m.upcomingTasks, contentWidth)
	upcomingBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.InfoColor).
		Padding(1, 2).
		Width(contentWidth).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			upcomingTitle,
			"",
			upcomingTasks,
		))

	// Combine all sections in the specified order
	var sections []string
	sections = append(sections, overviewBox)

	if len(m.highPriorityTasks) > 0 {
		sections = append(sections, highPrioBox)
	}

	if len(m.overdueTasks) > 0 {
		sections = append(sections, overdueBox)
	}

	if len(m.dueTodayTasks) > 0 {
		sections = append(sections, dueTodayBox)
	}

	if len(m.inProgressTasks) > 0 {
		sections = append(sections, inProgressBox)
	}

	if len(m.upcomingTasks) > 0 {
		sections = append(sections, upcomingBox)
	}

	// Place the dashboard in the center horizontally
	dashboardContent := lipgloss.JoinVertical(
		lipgloss.Left,
		sections...,
	)

	// Return content for viewport
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		mainBox.Render(dashboardContent),
	)
}

// Helper method to render tasks in a section
func (m *TodayDashboardModel) renderTasks(tasks []*models.Todo, width int) string {
	if len(tasks) == 0 {
		return styling.SubtextStyle.Render("No tasks")
	}

	var renderedTasks []string
	for _, task := range tasks {
		taskTitle := styling.TextStyle.Render(task.Title)

		// Add tags if present
		var tagStr string
		if len(task.Tags) > 0 {
			var tagRendered []string
			for _, tag := range task.Tags {
				tagRendered = append(tagRendered, styling.GetStyledTag(tag))
			}
			tagStr = lipgloss.JoinHorizontal(lipgloss.Left, tagRendered...)
		}

		leftContent := lipgloss.JoinHorizontal(
			lipgloss.Left,
			"◉ ",
			taskTitle,
		)
		rightContent := lipgloss.JoinHorizontal(
			lipgloss.Right,
			tagStr,
		)

		// Combine all elements
		taskLine := lipgloss.JoinHorizontal(
			lipgloss.Left,
			leftContent,
			lipgloss.NewStyle().Width(width-lipgloss.Width(leftContent)-6).Align(lipgloss.Right).Render(rightContent),
		)

		// Ensure the line fits within the width
		if lipgloss.Width(taskLine) > width-4 {
			taskLine = lipgloss.NewStyle().
				Width(width - 4).
				Render(taskLine)
		}

		renderedTasks = append(renderedTasks, taskLine)
	}

	return lipgloss.JoinVertical(lipgloss.Left, renderedTasks...)
}

// Loading view when data isn't ready
func (m *TodayDashboardModel) loadingView() string {
	content := lipgloss.NewStyle().
		Padding(1, 2).
		Foreground(theme.InfoColor).
		Bold(true).
		Render("Loading today's tasks...")

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
func (m *TodayDashboardModel) renderProgressBar() string {
	width := 20
	filledCount := int(float64(m.completedTasksCount) / float64(m.totalTasksCount) * float64(width))

	filled := strings.Repeat("▓", filledCount)
	empty := strings.Repeat("░", width-filledCount)

	percentage := 0
	if m.totalTasksCount > 0 {
		percentage = m.completedTasksCount * 100 / m.totalTasksCount
	}
	completedStats := m.translator.Tf("today_completed",
		map[string]interface{}{
			"completed": m.completedTasksCount,
			"total":     m.totalTasksCount,
			"percent":   percentage,
		})

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7DCFFF")).
		Render(filled) +
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89")).
			Render(empty) + "  " + completedStats

}

func (m *TodayDashboardModel) allTodosEmpty() bool {
	return len(m.highPriorityTasks) == 0 &&
		len(m.dueTodayTasks) == 0 &&
		len(m.inProgressTasks) == 0 &&
		len(m.overdueTasks) == 0 &&
		len(m.upcomingTasks) == 0
}

// ===========================================================================
// Messages
// ===========================================================================
type GetTodayDataMsg struct{}

type TodayDataUpdatedMsg struct{}

type GetCompletionStats struct{}

// ===========================================================================
// Commands
// ===========================================================================
func (m *TodayDashboardModel) GetCompletionStatsCmd() tea.Cmd {
	return func() tea.Msg {
		m.completedTasksCount, m.totalTasksCount, m.formattedTimeSpent = m.service.GetTodayCompletionStats()
		return TodayDataUpdatedMsg{}
	}
}

func (m *TodayDashboardModel) GetTodayDataCmd() tea.Cmd {
	return func() tea.Msg {
		var err error
		m.completedTasksCount, m.totalTasksCount, m.formattedTimeSpent = m.service.GetTodayCompletionStats()
		m.highPriorityTasks, m.dueTodayTasks, m.inProgressTasks, m.overdueTasks, m.upcomingTasks, err = m.service.GetTodosForToday()

		if err != nil {
			return TodoErrorMsg{err: err}
		}

		return TodayDataUpdatedMsg{}
	}
}

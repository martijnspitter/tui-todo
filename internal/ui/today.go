package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
)

type TodayDashboardModel struct {
	service    *service.AppService
	tuiService *service.TuiService
	translator *i18n.TranslationService
	width      int
	height     int

	// Section data
	highPriorityTasks []*models.Todo
	dueTodayTasks     []*models.Todo
	inProgressTasks   []*models.Todo
	upcomingTasks     []*models.Todo

	// Stats
	completedTasksCount int
	totalTasksCount     int

	// Currently selected section and item
	activeSection         int
	selectedItemInSection int
}

func NewTodayModel(service *service.AppService, tuiService *service.TuiService, translator *i18n.TranslationService) *TodayDashboardModel {
	return &TodayDashboardModel{
		service:    service,
		tuiService: tuiService,
		translator: translator,
	}
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

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7DCFFF")).
		Render(filled) +
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89")).
			Render(empty) +
		fmt.Sprintf("  %d/%d tasks complete (%d%%)",
			m.completedTasksCount,
			m.totalTasksCount,
			percentage)
}

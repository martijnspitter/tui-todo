package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

// Toast types
const (
	InfoToast    = "info"
	ErrorToast   = "error"
	WarningToast = "warning"
	SuccessToast = "success"
)

// Toast message
type ToastMsg struct {
	Message  string
	Type     string
	Duration time.Duration
}

// Message to dismiss toast
type DismissToastMsg struct{}

// Command to show toast
func showToastCmd(message, toastType string, duration time.Duration) tea.Cmd {
	return func() tea.Msg {
		return ToastMsg{
			Message:  message,
			Type:     toastType,
			Duration: duration,
		}
	}
}

// Show toast with default duration
func ShowDefaultToast(message, toastType string) tea.Cmd {
	return showToastCmd(message, toastType, 3*time.Second)
}

// Toast model
type ToastModel struct {
	Active  bool
	Message string
	Type    string
	Width   int
	Height  int
}

func NewToastModel() *ToastModel {
	return &ToastModel{}
}

func (m *ToastModel) Init() tea.Cmd {
	return nil
}

func (m *ToastModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ToastMsg:
		m.Message = msg.Message
		m.Type = msg.Type
		m.Active = true
		return m, tea.Sequence(
			tea.Tick(msg.Duration, func(time.Time) tea.Msg {
				return DismissToastMsg{}
			}),
		)

	case DismissToastMsg:
		m.Active = false

	case tea.WindowSizeMsg:
		// For toast, we don't need to store the full window size,
		// just ensure we have reasonable dimensions for the toast itself
		m.Width = max(msg.Width/3, 30)
		m.Height = 3
	}

	return m, nil
}

func (m *ToastModel) View() string {
	if !m.Active || m.Message == "" {
		return ""
	}

	var style lipgloss.Style

	// Define toast styles based on type
	switch m.Type {
	case ErrorToast:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(styling.ErrorColor).
			Padding(1, 1).
			Bold(true)

	case WarningToast:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(styling.WarningColor).
			Padding(1, 1).
			Bold(true)

	case SuccessToast:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(styling.SuccessColor).
			Padding(1, 1).
			Bold(true)

	default: // InfoToast
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(styling.InfoColor).
			Padding(1, 1).
			Bold(true)
	}

	return style.Render(m.Message)
}

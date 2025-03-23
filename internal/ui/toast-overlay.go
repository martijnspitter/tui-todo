package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

// ToastOverlay manages the toast model within an overlay
type ToastOverlay struct {
	toast   tea.Model
	overlay tea.Model
	content tea.Model
	width   int
	height  int
}

func NewToastOverlay(content tea.Model) *ToastOverlay {
	toastModel := NewToastModel()

	// Create the overlay but don't show it initially
	overlayModel := overlay.New(
		toastModel,
		content,
		overlay.Right,
		overlay.Bottom,
		-4,
		-2,
	)

	return &ToastOverlay{
		toast:   toastModel,
		overlay: overlayModel,
		content: content,
	}
}

func (m *ToastOverlay) Init() tea.Cmd {
	return nil
}

func (m *ToastOverlay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case ToastMsg, DismissToastMsg:
		m.toast, cmd = m.toast.Update(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.overlay, cmd = m.overlay.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ToastOverlay) View() string {
	return m.overlay.View()
}

package models

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
	active  bool
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
		active:  false,
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

	// First, update the toast model
	m.toast, cmd = m.toast.Update(msg)
	cmds = append(cmds, cmd)

	// Check if toast is active
	if toastModel, ok := m.toast.(*ToastModel); ok {
		m.active = toastModel.Active
	}

	// Then update the overlay with the same message
	m.overlay, cmd = m.overlay.Update(msg)
	cmds = append(cmds, cmd)

	// Update content
	m.content, cmd = m.content.Update(msg)
	cmds = append(cmds, cmd)

	// Handle window size for proper positioning
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = msg.Width
		m.height = msg.Height
	}

	// Return both commands
	return m, tea.Batch(cmds...)
}

func (m *ToastOverlay) View() string {
	return m.overlay.View()
}

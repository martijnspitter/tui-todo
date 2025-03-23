package styling

import "github.com/charmbracelet/lipgloss"

var (
	BrandColor     = lipgloss.Color("#cba6f7")
	HighlightColor = lipgloss.Color("#fab387")

	OpenStatusColor     = lipgloss.Color("#f5e0dc")
	DoingStatusColor    = lipgloss.Color("#89b4fa")
	DoneStatusColor     = lipgloss.Color("#a6e3a1")
	ArchivedStatusColor = lipgloss.Color("#9399b2")

	LowPriorityColor    = lipgloss.Color("#94e2d5")
	MediumPriorityColor = lipgloss.Color("#fab387")
	HighPriorityColor   = lipgloss.Color("#f38ba8")

	WhiteColor = lipgloss.Color("#fff")
	BlackColor = lipgloss.Color("#11111b")

	TextColor       = lipgloss.Color("#cdd6f4")
	SubtextColor    = lipgloss.Color("#a6adc8")
	BackgroundColor = lipgloss.Color("#313244")

	InfoColor    = lipgloss.Color("#186ddd")
	WarningColor = lipgloss.Color("#ff7c03")
	ErrorColor   = lipgloss.Color("#d13523")
	SuccessColor = lipgloss.Color("#1b7e41")

	FocusedStyle = lipgloss.NewStyle().Foreground(BrandColor)
	HoverStyle   = lipgloss.NewStyle().Foreground(HighlightColor)

	BorderWidth = 1
	Padding     = 1
)

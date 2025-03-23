package styling

import "github.com/charmbracelet/lipgloss"

var (
	BrandColorLight = lipgloss.Color("#2A9D8F")
	BrandColor      = lipgloss.Color("#3D7068")
	HighlightColor  = lipgloss.Color("#2A9D8F")

	TextColor            = lipgloss.Color("#EEF0F2")
	PlaceholderTextColor = lipgloss.Color("#888")

	InfoColor    = lipgloss.Color("#5386E4")
	WarningColor = lipgloss.Color("#FF7F11")
	ErrorColor   = lipgloss.Color("#FF1B1C")
	SuccessColor = lipgloss.Color("#3D7068")

	FocusedStyle = lipgloss.NewStyle().Foreground(HighlightColor)
	HoverStyle   = lipgloss.NewStyle().Foreground(BrandColorLight)

	BorderWidth = 1
	PaddingX    = 1
	PaddingY    = 1
)

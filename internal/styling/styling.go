package styling

import "github.com/charmbracelet/lipgloss"

var (
	BrandColorLight = lipgloss.Color("#7269b0")
	BrandColor      = lipgloss.Color("#3c327f")
	HighlightColor  = lipgloss.Color("#f7c76e")

	TextColor            = lipgloss.Color("#FFF")
	PlaceholderTextColor = lipgloss.Color("#888")

	InfoColor    = lipgloss.Color("#186ddd")
	WarningColor = lipgloss.Color("#ff7c03")
	ErrorColor   = lipgloss.Color("#d13523")
	SuccessColor = lipgloss.Color("#1b7e41")

	FocusedStyle = lipgloss.NewStyle().Foreground(HighlightColor)
	HoverStyle   = lipgloss.NewStyle().Foreground(BrandColorLight)

	BorderWidth = 1
	PaddingX    = 1
	PaddingY    = 1
)

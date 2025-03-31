package styling

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/models"
)

var (
	Mauve     = lipgloss.Color("#cba6f7")
	Yellow    = lipgloss.Color("#f9e2af")
	Lavender  = lipgloss.Color("#b4befe")
	Rosewater = lipgloss.Color("#f2cdcd")

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
	RowColor        = lipgloss.Color("#1e1e2e")

	InfoColor    = lipgloss.Color("#186ddd")
	WarningColor = lipgloss.Color("#ff7c03")
	ErrorColor   = lipgloss.Color("#d13523")
	SuccessColor = lipgloss.Color("#1b7e41")

	FocusedStyle = lipgloss.NewStyle().Foreground(Mauve)
	HoverStyle   = lipgloss.NewStyle().Foreground(Yellow)
	RowStyle     = lipgloss.NewStyle().Background(BlackColor)
	TextStyle    = lipgloss.NewStyle().Foreground(TextColor)
	SubtextStyle = lipgloss.NewStyle().Foreground(SubtextColor)

	BorderWidth = 1
	Padding     = 1
)

func GetStyledStatus(translatedStatus string, status models.Status, selected, omitNumber bool) string {
	statusColors := map[models.Status]lipgloss.Color{
		models.Open:  OpenStatusColor,
		models.Doing: DoingStatusColor,
		models.Done:  DoneStatusColor,
	}

	statusColor := statusColors[status]

	return GetStyledTagWithIndicator(int(status)+1, translatedStatus, statusColor, selected, omitNumber)
}

func GetStyledTagWithIndicator(num int, text string, color lipgloss.Color, selected, omitNumber bool) string {
	// Create indicator with number
	indicator := lipgloss.NewStyle().
		Foreground(BlackColor).
		Background(color).
		Padding(0, 1, 0, 0).
		Bold(true).
		Render(fmt.Sprintf("%d", num))

	// Text section (status name)
	var textStyle lipgloss.Style
	var leftCapStyle lipgloss.Style
	var rightCapStyle lipgloss.Style
	if selected {
		// Active tab
		textStyle = lipgloss.NewStyle().
			Foreground(BlackColor).
			Background(color).
			Padding(0, 0)
		leftCapStyle = lipgloss.NewStyle().
			Foreground(color).
			Padding(0, 0)
		rightCapStyle = lipgloss.NewStyle().
			Foreground(color).
			Padding(0, 0).
			MarginRight(2)
	} else {
		// Inactive tab
		textStyle = lipgloss.NewStyle().
			Foreground(SubtextColor).
			Background(BackgroundColor).
			Padding(0, 0)
		leftCapStyle = lipgloss.NewStyle().
			Foreground(color).
			Padding(0, 0)
		rightCapStyle = lipgloss.NewStyle().
			Foreground(BackgroundColor).
			Padding(0, 0).
			MarginRight(2)
	}

	if omitNumber {
		leftCapStyle.Foreground(BackgroundColor)
	}

	statusText := textStyle.Render(" " + text)
	leftCap := leftCapStyle.Render("")
	rightCap := rightCapStyle.Render("")

	if omitNumber {
		return lipgloss.JoinHorizontal(lipgloss.Center, leftCap, textStyle.Render(text), rightCap)
	}

	// Combine indicator and text
	return lipgloss.JoinHorizontal(lipgloss.Center, leftCap, indicator, statusText, rightCap)
}

func GetStyledPriority(translatedP string, p models.Priority, selected, hovered bool) string {
	priorityColors := []lipgloss.Color{
		LowPriorityColor,
		MediumPriorityColor,
		HighPriorityColor,
	}
	bgColor := BackgroundColor
	textColor := SubtextColor
	if selected {
		bgColor = priorityColors[p]
		textColor = BlackColor
	}
	if hovered {
		bgColor = Yellow
	}

	// Text section (status name)
	textStyle := lipgloss.NewStyle().
		Foreground(textColor).
		Background(bgColor).
		Width(8).
		Align(lipgloss.Center).
		MarginRight(1)

	return textStyle.Render(translatedP)
}

func GetStyledUpdatedAt(text string) string {
	textStyle := lipgloss.NewStyle().
		Foreground(Lavender).
		Background(BackgroundColor).
		Padding(0, 1).
		Align(lipgloss.Center).
		MarginRight(1)

	width := lipgloss.Width(text) + 2

	return textStyle.Width(width).Render(text)
}

func GetStyledDueDate(text string, priority models.Priority) string {
	priorityColors := []lipgloss.Color{
		LowPriorityColor,
		MediumPriorityColor,
		HighPriorityColor,
	}

	textStyle := lipgloss.NewStyle().
		Foreground(priorityColors[priority]).
		Background(BackgroundColor).
		Padding(0, 1).
		Align(lipgloss.Center).
		MarginRight(1)

	width := lipgloss.Width(text) + 2

	return textStyle.Width(width).Render(text)
}

func GetStyledTag(tag string) string {
	textStyle := lipgloss.NewStyle().
		Foreground(BlackColor).
		Background(Rosewater).
		Padding(0, 1).
		Align(lipgloss.Center).
		MarginRight(1)

	return textStyle.Render(tag)
}

func GetSelectedBlock(selected bool) string {
	if selected {
		return lipgloss.NewStyle().
			Foreground(Yellow).
			Background(Yellow).
			Padding(0, 1).
			Align(lipgloss.Center).
			MarginRight(1).
			Render("")
	}

	return lipgloss.NewStyle().
		Padding(0, 1).
		Align(lipgloss.Center).
		MarginRight(1).
		Render("")
}

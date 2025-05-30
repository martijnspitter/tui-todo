package styling

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/theme"
)

var (
	FocusedStyle = lipgloss.NewStyle().Foreground(theme.Mauve)
	HoverStyle   = lipgloss.NewStyle().Foreground(theme.Yellow)
	RowStyle     = lipgloss.NewStyle().Background(theme.BlackColor)
	TextStyle    = lipgloss.NewStyle().Foreground(theme.TextColor)
	SubtextStyle = lipgloss.NewStyle().Foreground(theme.SubtextColor)
	EmptyStyle   = lipgloss.NewStyle().Foreground(theme.Green)
	WarningStyle = lipgloss.NewStyle().Foreground(theme.ErrorColor)

	BorderWidth = 1
	Padding     = 1
)

func GetStyledStatus(translatedStatus string, status models.Status, selected, omitNumber, hovered bool) string {
	statusColor := status.Color()

	return GetStyledTagWithIndicator(int(status)+2, translatedStatus, statusColor, selected, omitNumber, hovered)
}

func GetStyledTagWithIndicator(num int, text string, color lipgloss.Color, selected, omitNumber, hovered bool) string {
	// Create indicator with number
	indicator := lipgloss.NewStyle().
		Foreground(theme.BlackColor).
		Background(color).
		Padding(0, 1, 0, 0).
		Bold(true).
		Render(fmt.Sprintf("%d", num))

	// Text section (status name)
	var textStyle lipgloss.Style
	var leftCapStyle lipgloss.Style
	var rightCapStyle lipgloss.Style
	if hovered {
		textStyle = lipgloss.NewStyle().
			Foreground(theme.BlackColor).
			Background(theme.Yellow).
			Padding(0, 0)
		leftCapStyle = lipgloss.NewStyle().
			Foreground(theme.Yellow).
			Padding(0, 0)
		rightCapStyle = lipgloss.NewStyle().
			Foreground(theme.Yellow).
			Padding(0, 0).
			MarginRight(2)
	} else if selected {
		textStyle = lipgloss.NewStyle().
			Foreground(theme.BlackColor).
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
			Foreground(theme.SubtextColor).
			Background(theme.BackgroundColor).
			Padding(0, 0)
		leftCapStyle = lipgloss.NewStyle().
			Foreground(color).
			Padding(0, 0)
		rightCapStyle = lipgloss.NewStyle().
			Foreground(theme.BackgroundColor).
			Padding(0, 0).
			MarginRight(2)
	}

	if omitNumber {
		leftCapStyle.Foreground(theme.BackgroundColor)
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
	bgColor := theme.BackgroundColor
	textColor := theme.SubtextColor
	if selected {
		bgColor = p.Color()
		textColor = theme.BlackColor
	}
	if hovered {
		bgColor = theme.Yellow
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
		Foreground(theme.Lavender).
		Background(theme.BackgroundColor).
		Padding(0, 1).
		Align(lipgloss.Center).
		MarginRight(1)

	width := lipgloss.Width(text) + 2

	return textStyle.Width(width).Render(text)
}

func GetStyledDueDate(text string, priority models.Priority) string {
	textStyle := lipgloss.NewStyle().
		Foreground(priority.Color()).
		Background(theme.BackgroundColor).
		Padding(0, 1).
		Align(lipgloss.Center).
		MarginRight(1)

	width := lipgloss.Width(text) + 2

	return textStyle.Width(width).Render(text)
}

func GetTimeSpend(text string) string {
	textStyle := lipgloss.NewStyle().
		Foreground(theme.Teal).
		Background(theme.BackgroundColor).
		Padding(0, 1).
		Align(lipgloss.Center)

	width := lipgloss.Width(text) + 2

	return textStyle.Width(width).Render(text)
}

func GetStyledTag(tag string) string {
	textStyle := lipgloss.NewStyle().
		Foreground(theme.BlackColor).
		Background(theme.Rosewater).
		Padding(0, 1).
		Align(lipgloss.Center).
		MarginRight(1)

	return textStyle.Render(tag)
}

func GetSelectedBlock(selected bool) string {
	if selected {
		return lipgloss.NewStyle().
			Foreground(theme.Yellow).
			Background(theme.Yellow).
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

func RenderMarkdown(md string) string {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	rendered, err := r.Render(md)
	if err != nil {
		// Fallback to raw markdown if rendering fails
		return md
	}

	return rendered

}

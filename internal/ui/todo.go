package ui

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type TodoItem struct {
	todo *models.Todo
}

func (i TodoItem) Title() string {
	return i.todo.Title
}

func (i TodoItem) Description() string {
	return i.todo.Description
}

func (i TodoItem) FilterValue() string {
	return i.todo.Title + " " + i.todo.Description
}

type TodoModel struct {
	translator *i18n.TranslationService
}

func (d TodoModel) Height() int                             { return 1 }
func (d TodoModel) Spacing() int                            { return 0 }
func (d TodoModel) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d TodoModel) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(TodoItem)
	if !ok {
		return
	}
	width := m.Width() - 4

	// Left-aligned elements
	selected := styling.GetSelectedBlock(index == m.Index())
	translatedPriority := d.translator.T(i.todo.Priority.String())
	priorityMarker := styling.GetStyledPriority(translatedPriority, i.todo.Priority, true, false)
	title := styling.TextStyle.MarginRight(1).Width(50).Render(truncateString(i.Title(), 50))

	leftElementsWidth := lipgloss.Width(selected) + lipgloss.Width(priorityMarker) + lipgloss.Width(title)

	if leftElementsWidth >= width {
		widthAvailableForTitle := width - lipgloss.Width(selected) - 1
		shortTitle := styling.TextStyle.MarginRight(1).Width(widthAvailableForTitle).Render(truncateString(i.Title(), widthAvailableForTitle))
		row := lipgloss.JoinHorizontal(lipgloss.Center, selected, shortTitle)
		fmt.Fprint(w, row)
		return
	}

	// Right-aligned elements
	var rightElements []string

	// Add tags
	tags := ""
	for _, tag := range i.todo.Tags {
		tags += styling.GetStyledTag(tag)
	}
	if tags != "" {
		rightElements = append(rightElements, tags)
	}

	// Add due date if present
	dueDate := ""
	if i.todo.DueDate != nil {
		translatedDueDate := d.translator.Tf("ui.due", map[string]interface{}{"Time": i.todo.DueDate.Format(time.Stamp)})
		dueDate = styling.GetStyledDueDate(translatedDueDate, i.todo.Priority)
		rightElements = append(rightElements, dueDate)
	}

	// Add updated at timestamp
	translatedUpdatedAt := d.translator.Tf("ui.updated", map[string]interface{}{"Time": i.todo.UpdatedAt.Format(time.Stamp)})
	updatedAt := styling.GetStyledUpdatedAt(translatedUpdatedAt)
	rightElements = append(rightElements, updatedAt)

	// Join right elements
	rightContent := lipgloss.JoinHorizontal(lipgloss.Right, rightElements...)
	rightWidth := lipgloss.Width(rightContent)

	// Calculate space for description
	descriptionMaxWidth := width - leftElementsWidth - rightWidth - 2 // 2 for some padding

	// Truncate description if needed
	description := i.todo.Description
	if descriptionMaxWidth > 20 {
		description = truncateString(description, descriptionMaxWidth)
	} else {
		description = ""
	}

	styledDescription := styling.SubtextStyle.Width(descriptionMaxWidth).Render(description)

	// Assemble the row with left content taking remaining space and right content aligned to the right
	leftContent := lipgloss.JoinHorizontal(lipgloss.Left, selected, priorityMarker, title, styledDescription)

	// Join everything, ensuring right alignment for the right content
	row := lipgloss.NewStyle().Width(width).Render(
		lipgloss.JoinHorizontal(lipgloss.Left,
			leftContent,
			lipgloss.NewStyle().Width(width-lipgloss.Width(leftContent)).Align(lipgloss.Right).Render(rightContent),
		),
	)

	fmt.Fprint(w, row)
}

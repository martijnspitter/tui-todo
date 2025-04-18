package ui

import (
	"fmt"
	"io"
	"strings"
	"time"

	"slices"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type TodoItem struct {
	todo       *models.Todo
	tuiService *service.TuiService
}

func (i *TodoItem) Title() string {
	return i.todo.Title
}

func (i *TodoItem) Description() string {
	// Replace newline characters with enter symbol
	cleanDesc := strings.ReplaceAll(i.todo.Description, "\n", " ↵ ")
	// Also handle Windows-style newlines
	cleanDesc = strings.ReplaceAll(cleanDesc, "\r", "")
	return cleanDesc
}

func (i *TodoItem) FilterValue() string {
	if i.tuiService.IsTagFilterActive() {
		return strings.Join(i.todo.Tags, " ")
	}
	return i.todo.Title + " " + i.todo.Description
}

type TodoModel struct {
	translator *i18n.TranslationService
	tuiService *service.TuiService
}

func (d TodoModel) Height() int                             { return 1 }
func (d TodoModel) Spacing() int                            { return 0 }
func (d TodoModel) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d TodoModel) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*TodoItem)
	if !ok {
		return
	}
	width := m.Width() - 4

	// Left-aligned elements
	selected := styling.GetSelectedBlock(index == m.Index())
	translatedPriority := d.translator.T(i.todo.Priority.String())
	priorityMarker := styling.GetStyledPriority(translatedPriority, i.todo.Priority, true, false)
	translatedStatus := d.translator.T(i.todo.Status.String())
	statusMarker := styling.GetStyledStatus(translatedStatus, i.todo.Status, true, true, false)
	if i.todo.Status != models.Doing {
		statusMarker = statusMarker + " "
	}
	// Add due date if present
	dueDate := ""
	if i.todo.DueDate != nil {
		translatedDueDate := d.translator.Tf("ui.due", map[string]interface{}{"Time": i.todo.DueDate.Format(time.Stamp)})
		dueDate = styling.GetStyledDueDate(translatedDueDate, i.todo.Priority)

	}

	statusLength := 0

	if d.tuiService.CurrentView == service.AllPane {
		statusLength = lipgloss.Width(statusMarker)
	}

	requiredItemsWidth := statusLength + lipgloss.Width(selected) + lipgloss.Width(priorityMarker)

	titleWidth, descriptionWidth, leftWidth, remainderWidth := d.tuiService.DetermineMaxWidthsForTodo(width, requiredItemsWidth, lipgloss.Width(dueDate))

	title := styling.TextStyle.MarginRight(1).Width(titleWidth).Render(truncateString(i.Title(), titleWidth))

	descStr := ""
	if descriptionWidth > 50 {
		descStr = styling.SubtextStyle.Width(descriptionWidth).Render(truncateString(i.Description(), descriptionWidth))
	} else {
		descStr = ""
	}

	if leftWidth >= width {
		widthAvailableForTitle := width - lipgloss.Width(selected) - 1
		shortTitle := styling.TextStyle.MarginRight(1).Width(widthAvailableForTitle).Render(truncateString(i.Title(), widthAvailableForTitle))
		row := lipgloss.JoinHorizontal(lipgloss.Center, selected, shortTitle)
		fmt.Fprint(w, row)
		return
	}

	var rightElements []string

	type elementInfo struct {
		element  string
		index    int
		priority int
	}
	var elementsToCheck []elementInfo

	tags := ""
	for _, tag := range i.todo.Tags {
		tags += styling.GetStyledTag(tag)
	}
	if tags != "" {
		elementsToCheck = append(elementsToCheck, struct {
			element  string
			index    int
			priority int
		}{tags, 2, 1})
	}

	elementsToCheck = append(elementsToCheck, struct {
		element  string
		index    int
		priority int
	}{dueDate, 1, 3})

	translatedUpdatedAt := d.translator.Tf("ui.updated", map[string]interface{}{"Time": i.todo.UpdatedAt.Format(time.Stamp)})
	updatedAt := styling.GetStyledUpdatedAt(translatedUpdatedAt)
	elementsToCheck = append(elementsToCheck, struct {
		element  string
		index    int
		priority int
	}{updatedAt, 4, 3})

	for {
		// Calculate total width of current elements
		totalWidth := 0
		for _, item := range elementsToCheck {
			totalWidth += lipgloss.Width(item.element)
		}

		// If everything fits, we're done
		if totalWidth <= remainderWidth || len(elementsToCheck) == 0 {
			break
		}

		// Find and remove lowest priority element
		lowestPriorityIdx := 0
		lowestPriority := 0

		for idx, item := range elementsToCheck {
			if item.index > lowestPriority {
				lowestPriority = item.index
				lowestPriorityIdx = idx
			}
		}

		// Remove the lowest priority element
		elementsToCheck = slices.Delete(elementsToCheck, lowestPriorityIdx, lowestPriorityIdx+1)
	}

	slices.SortFunc(elementsToCheck, func(a, b elementInfo) int {
		return a.priority - b.priority
	})

	for _, item := range elementsToCheck {
		rightElements = append(rightElements, item.element)
	}

	rightContent := lipgloss.JoinHorizontal(lipgloss.Right, rightElements...)
	leftContent := lipgloss.JoinHorizontal(lipgloss.Left, selected, priorityMarker, title, descStr)

	if d.tuiService.CurrentView == service.AllPane {
		leftContent = lipgloss.JoinHorizontal(lipgloss.Left, selected, priorityMarker, statusMarker, title, descStr)
	}

	row := lipgloss.NewStyle().Width(width).Render(
		lipgloss.JoinHorizontal(lipgloss.Left,
			leftContent,
			lipgloss.NewStyle().Width(width-lipgloss.Width(leftContent)).Align(lipgloss.Right).Render(rightContent),
		),
	)

	fmt.Fprint(w, row)
}

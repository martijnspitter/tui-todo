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
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type TagItem struct {
	tag *models.Tag
}

func (i *TagItem) Title() string {
	if i.tag == nil {
		return ""
	}
	return i.tag.Name
}

func (i *TagItem) Description() string {
	if i.tag == nil {
		return ""
	}
	return i.tag.Description
}

func (i *TagItem) FilterValue() string {
	return i.tag.Name
}

type TagModel struct {
	tuiService *service.TuiService
	translator *i18n.TranslationService
}

func (m TagModel) Height() int                             { return 1 }
func (m TagModel) Spacing() int                            { return 0 }
func (m TagModel) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (m TagModel) Render(w io.Writer, l list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*TagItem)
	if !ok {
		return
	}

	selected := styling.GetSelectedBlock(index == l.Index())
	translatedUpdatedAt := m.translator.Tf("ui.updated", map[string]interface{}{"Time": i.tag.UpdatedAt.Format(time.Stamp)})
	updatedAt := styling.GetStyledUpdatedAt(translatedUpdatedAt)
	requItemsWidth := lipgloss.Width(selected) + lipgloss.Width(updatedAt)
	nameWidth, descriptionWidth := m.tuiService.DetermineMaxWidthsForTag(l.Width()-4, requItemsWidth)
	name := styling.TextStyle.MarginRight(1).Width(nameWidth).Render(truncateString(i.Title(), nameWidth))
	description := styling.SubtextStyle.Width(descriptionWidth).Render(truncateString(i.Description(), descriptionWidth))

	leftContent := lipgloss.JoinHorizontal(lipgloss.Left, selected, name, description)
	rightContent := lipgloss.JoinHorizontal(lipgloss.Right, updatedAt)
	row := lipgloss.NewStyle().Width(l.Width() - 4).Render(
		lipgloss.JoinHorizontal(lipgloss.Left,
			leftContent,
			lipgloss.NewStyle().Width(l.Width()-4-lipgloss.Width(leftContent)).Align(lipgloss.Right).Render(rightContent),
		),
	)

	fmt.Fprintf(w, row)
}

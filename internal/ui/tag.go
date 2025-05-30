package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/martijnspitter/tui-todo/internal/models"
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
	return ""
}

func (i *TagItem) FilterValue() string {
	return i.tag.Name
}

type TagModel struct{}

func (m TagModel) Height() int                             { return 1 }
func (m TagModel) Spacing() int                            { return 0 }
func (m TagModel) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (m TagModel) Render(w io.Writer, _ list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*TagItem)
	if !ok {
		return
	}

	width := 20 // Fixed width for tags

	// Render the tag item
	title := i.tag.Name
	if len(title) > width {
		title = title[:width-3] + "..."
	}
	fmt.Fprintf(w, "%s", title)
}

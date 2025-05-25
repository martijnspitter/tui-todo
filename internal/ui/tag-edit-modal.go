package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
	"github.com/martijnspitter/tui-todo/internal/theme"
)

type tagState int

const (
	browsingTags tagState = iota
	creatingTag
	deletingTag
)

// TagItem represents a tag item in the list
type TagItem struct {
	name string
}

func (i TagItem) Title() string       { return i.name }
func (i TagItem) Description() string { return "" }
func (i TagItem) FilterValue() string { return i.name }

// TagManagementModal allows viewing and managing tags
type TagManagementModal struct {
	tagInput     textinput.Model
	tagList      list.Model
	state        tagState
	width        int
	height       int
	appService   *service.AppService
	tuiService   *service.TuiService
	translator   *i18n.TranslationService
	help         tea.Model
	selectedTag  string
	deletingTag  string
	confirmInput textinput.Model
}

func NewTagManagementModal(width, height int, appService *service.AppService, tuiService *service.TuiService, translationService *i18n.TranslationService) *TagManagementModal {
	help := NewHelpModel(appService, tuiService, translationService)

	// Setup tag input for creating new tags
	ti := textinput.New()
	ti.Placeholder = translationService.T("tag.new_placeholder")
	ti.CharLimit = 50

	// Setup confirmation input
	confirmInput := textinput.New()
	confirmInput.Placeholder = translationService.T("tag.confirm_delete")

	// Setup tag list
	tagItems := []list.Item{}

	// Create list model
	listModel := list.New(tagItems, list.NewDefaultDelegate(), width/2-12, height/2-10)
	listModel.Title = translationService.T("modal.manage_tags")
	listModel.SetShowHelp(false)

	return &TagManagementModal{
		tagInput:     ti,
		tagList:      listModel,
		state:        browsingTags,
		width:        width,
		height:       height,
		appService:   appService,
		tuiService:   tuiService,
		translator:   translationService,
		help:         help,
		confirmInput: confirmInput,
	}
}

func (m *TagManagementModal) Init() tea.Cmd {
	return tea.Batch(
		m.loadTags(),
		textinput.Blink,
	)
}

func (m *TagManagementModal) loadTags() tea.Cmd {
	return func() tea.Msg {
		tags, err := m.appService.GetAllTags()
		if err != nil {
			return TodoErrorMsg{err: err}
		}

		var items []list.Item
		for _, tag := range tags {
			items = append(items, TagItem{name: tag})
		}

		return TagsLoadedMsg{items: items}
	}
}

type TagsLoadedMsg struct {
	items []list.Item
}

type TagErrorMsg struct {
	err error
}

func (m *TagManagementModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case TagsLoadedMsg:
		m.tagList.SetItems(msg.items)

	case TagErrorMsg:
		// Handle error - maybe show a notification or toast
		log.Error("Tag error", "error", msg.err)

	case tea.KeyMsg:
		switch m.state {
		case browsingTags:
			switch {
			case key.Matches(msg, m.tuiService.KeyMap.Quit):
				// Close modal without saving
				return m, func() tea.Msg { return modalCloseMsg{reload: false} }

			case key.Matches(msg, m.tuiService.KeyMap.Create):
				// Enter create tag mode
				m.state = creatingTag
				m.tagInput.Focus()

			case key.Matches(msg, m.tuiService.KeyMap.Delete):
				// Enter delete mode if a tag is selected
				if i, ok := m.tagList.SelectedItem().(TagItem); ok {
					m.deletingTag = i.name
					m.state = deletingTag
					m.confirmInput.SetValue("")
					m.confirmInput.Focus()
				}
			}

			// Let the list handle input in browse mode
			var cmd tea.Cmd
			m.tagList, cmd = m.tagList.Update(msg)
			cmds = append(cmds, cmd)

		case creatingTag:
			switch {
			case key.Matches(msg, m.tuiService.KeyMap.Quit):
				// Cancel tag creation
				m.tagInput.SetValue("")
				m.tagInput.Blur()
				m.state = browsingTags

			case key.Matches(msg, m.tuiService.KeyMap.Save):
				// Save the new tag
				if strings.TrimSpace(m.tagInput.Value()) != "" {
					cmds = append(cmds, m.createTag(m.tagInput.Value()))
				}
				m.tagInput.SetValue("")
				m.tagInput.Blur()
				m.state = browsingTags
			}

			// Update the text input
			var cmd tea.Cmd
			m.tagInput, cmd = m.tagInput.Update(msg)
			cmds = append(cmds, cmd)

		case deletingTag:
			switch {
			case key.Matches(msg, m.tuiService.KeyMap.Quit):
				// Cancel deletion
				m.confirmInput.SetValue("")
				m.confirmInput.Blur()
				m.state = browsingTags

			case key.Matches(msg, m.tuiService.KeyMap.Delete):
				if strings.ToLower(strings.TrimSpace(m.confirmInput.Value())) == "delete" {
					// Confirm deletion
					cmds = append(cmds, m.deleteTag(m.deletingTag))
				}
				m.confirmInput.SetValue("")
				m.confirmInput.Blur()
				m.state = browsingTags
			}

			// Update the confirm input
			var cmd tea.Cmd
			m.confirmInput, cmd = m.confirmInput.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.tagList.SetWidth(msg.Width/2 - 12)
		m.tagList.SetHeight(msg.Height/2 - 10)
	}

	return m, tea.Batch(cmds...)
}

func (m *TagManagementModal) createTag(tagName string) tea.Cmd {
	return func() tea.Msg {
		err := m.appService.CreateTag(tagName)
		if err != nil {
			return TagErrorMsg{err: err}
		}

		// Reload tags after creation
		return m.loadTags()()
	}
}

func (m *TagManagementModal) deleteTag(id int64) tea.Cmd {
	return func() tea.Msg {
		err := m.appService.DeleteTag(id)
		if err != nil {
			return TagErrorMsg{err: err}
		}

		// Reload tags after deletion
		return m.loadTags()()
	}
}

func (m *TagManagementModal) View() string {
	// Create modal style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(m.width / 2).
		BorderForeground(theme.Mauve)

	var content string

	switch m.state {
	case browsingTags:
		content = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			m.tagList.View(),
			styling.TextStyle.Render(m.translator.T("tag.manage_hint")),
			m.help.View(),
		)

	case creatingTag:
		content = fmt.Sprintf(
			"%s\n\n%s\n%s\n\n%s",
			styling.FocusedStyle.Render(m.translator.T("tag.create_new")),
			m.tagInput.View(),
			styling.TextStyle.Render(m.translator.T("tag.create_hint")),
			m.help.View(),
		)

	case deletingTag:
		content = fmt.Sprintf(
			"%s\n\n%s\n%s\n\n%s\n\n%s",
			styling.WarningStyle.Render(m.translator.T("tag.delete_warning")),
			m.translator.Tf("tag.delete_confirm", map[string]interface{}{"TagName": m.deletingTag}),
			styling.TextStyle.Render(m.translator.T("tag.type_delete")),
			m.confirmInput.View(),
			m.help.View(),
		)
	}

	// Center the modal
	positioned := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(content),
	)

	return positioned
}

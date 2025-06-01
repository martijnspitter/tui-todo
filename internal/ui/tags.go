package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/styling"
)

type TagsModel struct {
	service    *service.AppService
	tuiService *service.TuiService
	translator *i18n.TranslationService
	list       list.Model
	width      int
	height     int
}

func NewTagsModel(service *service.AppService, tuiService *service.TuiService, translator *i18n.TranslationService) *TagsModel {
	// Setup list
	tagList := list.New([]list.Item{}, TagModel{tuiService, translator}, 0, 0)
	tagList.Title = ""
	tagList.DisableQuitKeybindings()
	tagList.SetShowTitle(false)
	tagList.SetShowHelp(false)
	tagList.SetShowStatusBar(false)
	tagList.SetFilteringEnabled(true)

	return &TagsModel{
		service:    service,
		tuiService: tuiService,
		translator: translator,
		list:       tagList,
	}
}

func (m *TagsModel) Init() tea.Cmd {
	return nil
}

func (m *TagsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.tuiService.KeyMap.Edit):
			// Edit selected Tag
			if m.shouldAllowTagCrud() {
				item := m.list.SelectedItem().(*TagItem)
				return m, m.showEditModalCmd(item.tag)
			}
		case key.Matches(msg, m.tuiService.KeyMap.Delete):
			// Delete selected todo
			if m.shouldAllowTagCrud() {
				item := m.list.SelectedItem().(*TagItem)
				return m, m.showConfirmDeleteCmd(item.tag.ID)
			}
		case key.Matches(msg, m.tuiService.KeyMap.New):
			// Create new Todo
			if m.tuiService.CurrentView == service.TagsPane {
				tag := &models.Tag{ID: -1}
				return m, m.showEditModalCmd(tag)
			}
		}
	case RemoveFilterMsg:
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	case tagsLoadedMsg:
		// Update the list with new tags
		items := make([]list.Item, len(msg.tags))
		for i, tag := range msg.tags {
			items[i] = &TagItem{tag: tag}
		}
		cmd := m.list.SetItems(items)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		headerHeight := 3 // Title + top border
		footerHeight := 3 // Input + bottom padding
		m.width = msg.Width
		m.height = msg.Height - headerHeight - footerHeight

		m.list.SetSize(msg.Width, m.height)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *TagsModel) View() string {
	listView := lipgloss.NewStyle().Width(m.width - 2).Padding(styling.Padding).Render(m.list.View())
	if len(m.list.Items()) == 0 {
		listView = EmptyNothingFoundView(m.translator, m.width, m.height)
	}

	return listView
}

// ===========================================================================
// Helpers
// ===========================================================================
func (m *TagsModel) shouldAllowTagCrud() bool {
	return m.list.SelectedItem() != nil && m.tuiService.CurrentView == service.TagsPane
}

func (m *TagsModel) SetHeight(height int) {
	m.height = height
	m.list.SetHeight(height)
}

// ===========================================================================
// Commands
// ===========================================================================
func (m *TagsModel) showEditModalCmd(tag *models.Tag) tea.Cmd {
	return func() tea.Msg {
		m.tuiService.SwitchToEditTagView()
		modalComponent := NewTagEditModal(tag, m.width, m.height, m.service, m.tuiService, m.translator)
		return showModalMsg{
			modal: modalComponent,
		}
	}
}

func (m *TagsModel) showConfirmDeleteCmd(todoID int64) tea.Cmd {
	return func() tea.Msg {
		m.tuiService.SwitchToConfirmDeleteView()
		modalComponent := NewConfirmDeleteModal(m.service, m.tuiService, m.translator, todoID, true)
		return showModalMsg{
			modal: modalComponent,
		}
	}
}

package ui

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/keys"
	"github.com/martijnspitter/tui-todo/internal/service"
)

type HelpModel struct {
	service    *service.AppService
	tuiService *service.TuiService
	translator *i18n.TranslationService
	help       help.Model
	width      int
}

func NewHelpModel(service *service.AppService, tuiService *service.TuiService, translator *i18n.TranslationService) *HelpModel {
	return &HelpModel{
		help:       help.New(),
		service:    service,
		tuiService: tuiService,
		translator: translator,
	}
}

func (m *HelpModel) Init() tea.Cmd {
	return nil
}

func (m *HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

func (m *HelpModel) View() string {
	keyMap := m.getContextualKeyMap()

	// Use the help model to render them
	return m.help.View(keyMap)
}

func (m *HelpModel) getContextualKeyMap() keys.HelpKeyMap {
	currentView := m.tuiService.CurrentView
	filterState := m.tuiService.FilterState
	baseKeyMap := m.tuiService.KeyMap

	// Create a new context-specific key map
	contextKeyMap := keys.NewHelpKeyMap(m.translator)

	// Always show these keys regardless of context when not filtering
	if !filterState.IsFilterActive && currentView != service.AddEditModal {
		contextKeyMap.AddBindingInShort(baseKeyMap.Help)
		contextKeyMap.AddBindingInShort(baseKeyMap.Quit)
	}

	// Add view-specific bindings
	switch currentView {
	case service.OpenPane, service.DoingPane, service.DonePane, service.AllPane:
		if filterState.IsFilterActive {
			contextKeyMap.AddBindingInShort(baseKeyMap.Cancel)
		} else {
			// List view shows navigation keys
			contextKeyMap.AddBindingInShort(baseKeyMap.New)
			contextKeyMap.AddBindingInShort(baseKeyMap.Filter)

			contextKeyMap.AddBindingInFull(baseKeyMap.Up)
			contextKeyMap.AddBindingInFull(baseKeyMap.Down)
			contextKeyMap.AddBindingInFull(baseKeyMap.Filter)
			contextKeyMap.AddBindingInFull(baseKeyMap.Help)
			contextKeyMap.AddBindingInFull(baseKeyMap.Home)
			contextKeyMap.AddBindingInFull(baseKeyMap.End)

			contextKeyMap.AddBindingInFull(baseKeyMap.SwitchPane)
			contextKeyMap.AddBindingInFull(baseKeyMap.New)
			contextKeyMap.AddBindingInFull(baseKeyMap.Edit)

			// Show these in full help or when relevant
			contextKeyMap.AddBindingInFull(baseKeyMap.Delete)
			contextKeyMap.AddBindingInFull(baseKeyMap.AdvanceStatus)
			contextKeyMap.AddBindingInFull(baseKeyMap.Archive)
		}

		// Only show archived toggle in All filter mode
		if currentView == service.AllPane {
			contextKeyMap.AddBindingInShort(baseKeyMap.ToggleArchived)
			contextKeyMap.AddBindingInFull(baseKeyMap.ToggleArchived)
		}

	case service.AddEditModal:
		// Edit view shows edit-specific keys
		contextKeyMap.AddBindingInShort(baseKeyMap.Cancel)
		contextKeyMap.AddBindingInShort(baseKeyMap.Next)
		contextKeyMap.AddBindingInShort(baseKeyMap.Prev)
		contextKeyMap.AddBindingInShort(baseKeyMap.Select)
		contextKeyMap.AddBindingInShort(baseKeyMap.Save)

	case service.ConfirmDeleteModal:
		// Confirm delete shows minimal keys
		contextKeyMap.AddBindingInFull(baseKeyMap.Select) // Confirm
	}

	return contextKeyMap
}

func (m *HelpModel) ToggleShowAll() {
	m.help.ShowAll = !m.help.ShowAll
}

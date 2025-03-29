package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Next           key.Binding
	Prev           key.Binding
	SwitchPane     key.Binding
	Select         key.Binding
	Quit           key.Binding
	New            key.Binding
	Edit           key.Binding
	Delete         key.Binding
	AdvanceStatus  key.Binding
	Archive        key.Binding
	ToggleArchived key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Next: key.NewBinding(
			key.WithKeys("right", "tab"),
			key.WithHelp("right/tab", "Go next"),
		),
		Prev: key.NewBinding(
			key.WithKeys("left", "shift+tab"),
			key.WithHelp("left/shift+tab", "Go previous"),
		),
		SwitchPane: key.NewBinding(
			key.WithKeys("1", "2", "3", "4", "5", "6"),
			key.WithHelp("1-6", "Select Pane"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "esc"),
			key.WithHelp("ctrl+c/esc", "Close application"),
		),
		New: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "New Todo"),
		),
		Edit: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "Edit Todo"),
		),
		Delete: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "Delete Todo"),
		),
		AdvanceStatus: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "Save"),
		),
		Archive: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "archive/unarchive"),
		),
		ToggleArchived: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "toggle archived todos"),
		),
	}
}

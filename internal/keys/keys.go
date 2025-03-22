package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Next        key.Binding
	Prev        key.Binding
	SwitchPane  key.Binding
	Select      key.Binding
	Impersonate key.Binding
	Quit        key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Next: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("right/l", "Go next"),
		),
		Prev: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("left/h", "Go previous"),
		),
		SwitchPane: key.NewBinding(
			key.WithKeys("1", "2", "3", "4", "5", "6"),
			key.WithHelp("1-6", "Select Pane"),
		),
		Impersonate: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("ctrl+i", "Impersonate"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "esc"),
			key.WithHelp("ctrl+c/esc", "Close application"),
		),
	}
}

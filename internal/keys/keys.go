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
	Help           key.Binding
	Filter         key.Binding
	Up             key.Binding
	Down           key.Binding
	Cancel         key.Binding
	Home           key.Binding
	End            key.Binding
	TagFilter      key.Binding
	Save           key.Binding
	About          key.Binding
	PageDown       key.Binding
	PageUp         key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Next: key.NewBinding(
			key.WithKeys("right", "tab"),
			key.WithHelp("right/tab", "help.right_tab"),
		),
		Prev: key.NewBinding(
			key.WithKeys("left", "shift+tab"),
			key.WithHelp("left/shift+tab", "help.left_shift_tab"),
		),
		SwitchPane: key.NewBinding(
			key.WithKeys("1", "2", "3", "4", "5"),
			key.WithHelp("1-5", "help.pane"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "help.enter"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "esc"),
			key.WithHelp("ctrl+c/esc", "help.ctrl_c_esc"),
		),
		New: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "help.ctrl_n"),
		),
		Edit: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "help.ctrl_e"),
		),
		Delete: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "help.ctrl_d"),
		),
		AdvanceStatus: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "help.ctrl_s"),
		),
		Archive: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "help.ctrl_a"),
		),
		ToggleArchived: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "help.a"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help.toggle"),
		),
		Up: key.NewBinding(
			key.WithKeys("k", "↑"),
			key.WithHelp("↑/k", "help.up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "↓"),
			key.WithHelp("↓/j", "help.down"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", " ", "f"),
			key.WithHelp("f/pgdn", "help.page_down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b"),
			key.WithHelp("b/pgup", "help.page_up"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "help.filter"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "help.cancel"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "help.home"),
		),
		End: key.NewBinding(
			key.WithKeys("G", "end"),
			key.WithHelp("G/end", "help.end"),
		),
		TagFilter: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "help.tag_filter"),
		),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "help.save"),
		),
		About: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "help.i"),
		),
	}
}

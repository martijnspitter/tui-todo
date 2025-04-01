package keys

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/martijnspitter/tui-todo/internal/i18n"
)

// HelpKeyMap is a key map that supports translation and contextual display
type HelpKeyMap struct {
	shortBindings []key.Binding
	fullBindings  [][]key.Binding
	translator    *i18n.TranslationService
}

func NewHelpKeyMap(translator *i18n.TranslationService) HelpKeyMap {
	return HelpKeyMap{
		shortBindings: []key.Binding{},
		fullBindings:  [][]key.Binding{{}},
		translator:    translator,
	}
}

// AddBindingInFull adds a key binding to the help keymap with translation support
func (k *HelpKeyMap) AddBindingInFull(binding key.Binding) {
	// Translate the help text if possible
	if k.translator != nil {
		translated := k.translator.T(binding.Help().Desc)
		binding.SetHelp(binding.Help().Key, translated)
	}

	if len(k.fullBindings) == 0 {
		k.fullBindings = append(k.fullBindings, []key.Binding{})
	}
	k.fullBindings[0] = append(k.fullBindings[0], binding)
}

func (k *HelpKeyMap) AddBindingInShort(binding key.Binding) {
	// Translate the help text if possible
	if k.translator != nil {
		translated := k.translator.T(binding.Help().Desc)
		binding.SetHelp(binding.Help().Key, translated)
	}

	k.shortBindings = append(k.shortBindings, binding)
}

// ShortHelp is what's displayed when help is first requested
func (k HelpKeyMap) ShortHelp() []key.Binding {
	return k.shortBindings
}

// FullHelp is the taller help view listing all key bindings
func (k HelpKeyMap) FullHelp() [][]key.Binding {
	return k.fullBindings
}

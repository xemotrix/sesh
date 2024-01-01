package creator

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit    key.Binding
	Confirm key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Confirm, k.Quit},
	}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
}

package app

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	NextPane key.Binding
	PrevPane key.Binding
	Down     key.Binding
	Up       key.Binding
	Insert   key.Binding
	Normal   key.Binding
	Run      key.Binding
	Save     key.Binding
	Open     key.Binding
	Help     key.Binding
	Quit     key.Binding
}

var DefaultKeyMap = KeyMap{
	NextPane: key.NewBinding(
		key.WithKeys("tab", "l"),
		key.WithHelp("tab/l", "next pane"),
	),
	PrevPane: key.NewBinding(
		key.WithKeys("shift+tab", "h"),
		key.WithHelp("shift+tab/h", "prev pane"),
	),
	Down: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "down"),
	),
	Up: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "up"),
	),
	Insert: key.NewBinding(
		key.WithKeys("i", "enter"),
		key.WithHelp("i/enter", "insert mode"),
	),
	Normal: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "normal mode"),
	),
	Run: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "run request"),
	),
	Save: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "save"),
	),
	Open: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open saved"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

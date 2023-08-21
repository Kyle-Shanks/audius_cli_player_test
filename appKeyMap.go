package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	// Global keys
	Help key.Binding
	Quit key.Binding

	// View keys
	Trending key.Binding
	Search   key.Binding

	// Table keys
	Up   key.Binding
	Down key.Binding
	// Left   key.Binding
	// Right  key.Binding
	Top    key.Binding
	Bottom key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Top, k.Bottom},
		{k.Help, k.Quit},
	}
}

var AppKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	// Left: key.NewBinding(
	// 	key.WithKeys("left", "h"),
	// 	key.WithHelp("←/h", "move left"),
	// ),
	// Right: key.NewBinding(
	// 	key.WithKeys("right", "l"),
	// 	key.WithHelp("→/l", "move right"),
	// ),
	Top: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "jump to top"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "jump to bottom"),
	),

	Trending: key.NewBinding(
		key.WithKeys("T"),
		key.WithHelp("T", "trending"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		// key.WithKeys("esc", "q", "ctrl+c"),
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "quit"),
	),
}

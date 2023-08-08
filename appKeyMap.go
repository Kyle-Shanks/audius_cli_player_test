package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	// Global keys
	Help key.Binding
	Quit key.Binding

	// Player keys
	Play key.Binding
	Mute key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Play, k.Mute},
		{k.Help, k.Quit},
	}
}

var AppKeyMap = KeyMap{
	// Up: key.NewBinding(
	// 	key.WithKeys("up", "k"),
	// 	key.WithHelp("↑/k", "move up"),
	// ),
	// Down: key.NewBinding(
	// 	key.WithKeys("down", "j"),
	// 	key.WithHelp("↓/j", "move down"),
	// ),
	// Left: key.NewBinding(
	// 	key.WithKeys("left", "h"),
	// 	key.WithHelp("←/h", "move left"),
	// ),
	// Right: key.NewBinding(
	// 	key.WithKeys("right", "l"),
	// 	key.WithHelp("→/l", "move right"),
	// ),
	Play: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "play/pause"),
	),
	Mute: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "toggle mute"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

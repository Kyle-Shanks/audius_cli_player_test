package player

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	// For handling things on app quit
	Quit key.Binding

	Pause       key.Binding
	Mute        key.Binding
	Repeat      key.Binding
	VolumeUp    key.Binding
	VolumeDown  key.Binding
	SkipForward key.Binding
	SkipBack    key.Binding
	Next        key.Binding
	Prev        key.Binding

	// Just for help text
	Play key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Play, k.Pause, k.Repeat, k.Mute},
		// {k.SkipBack, k.SkipForward},
	}
}

// TODO: Add Next and Prev keys
var PlayerKeyMap = KeyMap{
	Play: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "play track"),
	),
	Pause: key.NewBinding(
		key.WithKeys("p", " "),
		key.WithHelp("p/space", "toggle pause"),
	),
	SkipForward: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "10s forward"),
	),
	SkipBack: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "10s back"),
	),
	Mute: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "toggle mute"),
	),
	Repeat: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "toggle repeat"),
	),
	VolumeUp: key.NewBinding(
		key.WithKeys("."),
		key.WithHelp(".", "volume up"),
	),
	VolumeDown: key.NewBinding(
		key.WithKeys(","),
		key.WithHelp(",", "volume down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "q", "ctrl+c"),
		key.WithHelp("esc", "quit"),
	),
}

package common

import tea "github.com/charmbracelet/bubbletea"

type PlayTrackMsg struct {
	Track Track
}

func PlayTrackCmd(track Track) tea.Cmd {
	return func() tea.Msg {
		return PlayTrackMsg{
			Track: track,
		}
	}
}

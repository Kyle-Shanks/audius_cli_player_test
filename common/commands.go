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

type PlayTracksMsg struct {
	Tracks   []Track
	QueuePos int
}

func PlayTracksCmd(tracks []Track, pos int) tea.Cmd {
	return func() tea.Msg {
		return PlayTracksMsg{
			Tracks:   tracks,
			QueuePos: pos,
		}
	}
}

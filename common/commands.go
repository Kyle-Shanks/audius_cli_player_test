package common

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

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

func LogCmd(log string) tea.Cmd {
	return func() tea.Msg {
		f, err := tea.LogToFile(filepath.Join(GetDataPath(), "debug.log"), "debug")

		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}

		f.WriteString(log + "\n---\n")
		defer f.Close()

		return nil
	}
}

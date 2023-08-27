package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

type TracksTableView struct {
	focused     bool
	tracksTable TracksTable
	viewTitle   string
	fetchTracks func() ([]Track, error)
}

type tracksResponseMsg struct {
	viewTitle string
	tracks    []Track
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

func (tt *TracksTableView) fetchTracksCmd() tea.Cmd {
	return func() tea.Msg {
		tracks, err := tt.fetchTracks()

		if err != nil {
			return errMsg{err}
		}

		return tracksResponseMsg{
			viewTitle: tt.viewTitle,
			tracks:    tracks,
		}
	}
}

func NewTracksTableView(viewTitle string, fetchTracks func() ([]Track, error)) TracksTableView {
	t := NewTracksTable(
		13,
		func(track Track) tea.Cmd {
			return nil
		},
	)
	t.SetIsLoading(true)
	t.Focus()

	tt := TracksTableView{
		viewTitle:   viewTitle,
		tracksTable: t,
		focused:     false,
		fetchTracks: fetchTracks,
	}

	return tt
}

func (tt TracksTableView) Init() tea.Cmd {
	return tt.fetchTracksCmd()
}

func (tt TracksTableView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tracksResponseMsg:
		if msg.viewTitle == tt.viewTitle {
			tt.tracksTable.UpdateTracks(msg.tracks)
			tt.tracksTable.SetIsLoading(false)
		}
	}

	if tt.focused {
		res, cmd := tt.tracksTable.Update(msg)
		tt.tracksTable = res.(TracksTable)
		cmds = append(cmds, cmd)
	}

	return tt, tea.Batch(cmds...)
}

func (tt TracksTableView) View() string {
	return tt.tracksTable.View()
}

func (tt *TracksTableView) Focus() {
	tt.focused = true
}

func (tt *TracksTableView) Blur() {
	tt.focused = false
}

func (tt *TracksTableView) Focused() bool {
	return tt.focused
}

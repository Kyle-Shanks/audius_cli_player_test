package home

import (
	"app1/api"
	"app1/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UndergroundTracks struct {
	focused           bool
	tracksTable       common.TracksTable
	undergroundTracks []common.Track
}

type undergroundTracksResponseMsg struct {
	tracks []common.Track
}

func fetchUndergroundTracksCmd() tea.Msg {
	tracks, err := api.GetUndergroundTracks()

	if err != nil {
		return errMsg{err}
	}

	return undergroundTracksResponseMsg{
		tracks: tracks,
	}
}

func NewUndergroundTracks() UndergroundTracks {
	t := common.NewTracksTable(
		13,
		func(track common.Track) tea.Cmd {
			return nil
		},
	)
	t.SetIsLoading(true)
	t.Focus()

	tt := UndergroundTracks{tracksTable: t, focused: false}

	return tt
}

func (tt UndergroundTracks) Init() tea.Cmd {
	return fetchUndergroundTracksCmd
}

func (tt UndergroundTracks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case undergroundTracksResponseMsg:
		tt.tracksTable.UpdateTracks(msg.tracks)
		tt.tracksTable.SetIsLoading(false)
	}

	if tt.focused {
		res, cmd := tt.tracksTable.Update(msg)
		tt.tracksTable = res.(common.TracksTable)
		cmds = append(cmds, cmd)
	}

	return tt, tea.Batch(cmds...)
}

func (tt UndergroundTracks) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		// common.Header().Render("Underground Tracks"),
		tt.tracksTable.View(),
	)
}

func (tt *UndergroundTracks) Focus() {
	tt.focused = true
}

func (tt *UndergroundTracks) Blur() {
	tt.focused = false
}

func (tt *UndergroundTracks) Focused() bool {
	return tt.focused
}

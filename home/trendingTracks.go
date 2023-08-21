package home

import (
	"app1/api"
	"app1/common"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TrendingTracks struct {
	focused        bool
	tracksTable    common.TracksTable
	trendingTracks []common.Track
}

type trendingTracksResponseMsg struct {
	tracks []common.Track
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

func fetchTrendingTracksCmd() tea.Msg {
	tracks, err := api.GetTrendingTracks()

	if err != nil {
		return errMsg{err}
	}

	return trendingTracksResponseMsg{
		tracks: tracks,
	}
}

func NewTrendingTracks() TrendingTracks {
	t := common.NewTracksTable(
		13,
		func(track common.Track) tea.Cmd {
			return nil
		},
	)
	t.SetIsLoading(true)
	t.Focus()

	tt := TrendingTracks{tracksTable: t, focused: false}

	return tt
}

func (tt TrendingTracks) Init() tea.Cmd {
	return fetchTrendingTracksCmd
}

func (tt TrendingTracks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case trendingTracksResponseMsg:
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

func (tt TrendingTracks) View() string {
	header := common.ActiveHeader().Render("Trending Tracks")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		tt.tracksTable.View(),
	)
}

func (tt *TrendingTracks) Focus() {
	tt.focused = true
}

func (tt *TrendingTracks) Blur() {
	tt.focused = false
}

func (tt *TrendingTracks) Focused() bool {
	return tt.focused
}

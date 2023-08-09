package home

import (
	"app1/api"
	"app1/common"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TrendingTracks struct {
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

func (tt TrendingTracks) Init() tea.Cmd {
	tt.tracksTable.SetIsLoading(true)
	return fetchTrendingTracksCmd
}

func (tt TrendingTracks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case trendingTracksResponseMsg:
		tt.tracksTable.UpdateTracks(msg.tracks)
		tt.tracksTable.SetIsLoading(false)
	}

	res, cmd := tt.tracksTable.Update(msg)
	tt.tracksTable = res.(common.TracksTable)
	cmds = append(cmds, cmd)

	return tt, tea.Batch(cmds...)
}

func (tt TrendingTracks) View() string {
	header := lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true).
		Padding(0, 2).
		Render("Trending Tracks")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		tt.tracksTable.View(),
	)
}

func NewTrendingTracks() TrendingTracks {
	rows := []table.Row{
		{"", "Loading...", "", ""},
	}

	t := common.NewTracksTable(
		common.DefaultTracksTableColumns,
		rows,
		func(track common.Track) tea.Cmd {
			return nil
		},
	)

	tt := TrendingTracks{tracksTable: t}

	return tt
}

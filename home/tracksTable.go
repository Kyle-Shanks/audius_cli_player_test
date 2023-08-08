package home

import (
	"app1/api"
	"app1/common"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	tableContainerStyle         = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	tableInactiveContainerStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("243"))
)

type TracksTable struct {
	focused        bool
	isLoading      bool
	table          table.Model
	trendingTracks []common.Track
}

func (tt TracksTable) Focus() {
	tt.focused = true
}

func (tt TracksTable) Blur() {
	tt.focused = false
}

func (tt TracksTable) Focused() bool {
	return tt.focused
}

func (tt TracksTable) SetIsLoading(val bool) {
	tt.isLoading = val
}

func (tt TracksTable) IsLoading() bool {
	return tt.isLoading
}

type trendingTracksResponseMsg struct {
	tracks []common.Track
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

func fetchTrendingTracks() tea.Msg {
	tracks, err := api.GetTrendingTracks()

	if err != nil {
		return errMsg{err}
	}

	return trendingTracksResponseMsg{
		tracks: tracks,
	}
}

func (tt TracksTable) Init() tea.Cmd {
	return fetchTrendingTracks
}

func (tt TracksTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// tt.resizeTable(msg.Width)
		return tt, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			// row := tt.table.SelectedRow()
			track := tt.trendingTracks[tt.table.Cursor()]
			// trackId := tt.trendingTrackIds[tt.table.Cursor()]
			// fmt.Println(trackId)
			cmds = append(cmds, common.PlayTrackCmd(track))
		}

	case trendingTracksResponseMsg:
		tt.updateTableContents(msg.tracks)
	}

	tt.table, cmd = tt.table.Update(msg)
	cmds = append(cmds, cmd)

	return tt, tea.Batch(cmds...)
}

func (tt *TracksTable) updateTableContents(tracks []common.Track) {
	var trackRows []table.Row
	// var trackIds []string

	for i, track := range tracks {
		// trackIds = append(trackIds, track.Id)
		trackRows = append(
			trackRows,
			table.Row{
				fmt.Sprint(i + 1),
				common.RemoveEmojis(
					track.Title,
				),
				common.RemoveEmojis(
					track.User.Name,
				),
				common.GetLengthText(track.Duration),
			},
		)
	}

	tt.trendingTracks = tracks
	// tt.trendingTrackIds = trackIds
	tt.table.SetRows(trackRows)
}

func (tt TracksTable) View() string {
	if tt.Focused() {
		return tableContainerStyle.Render(tt.table.View())
	} else {
		return tableInactiveContainerStyle.Render(tt.table.View())
	}
}

func (tt TracksTable) resizeTable(width int) {
	// TODO: Figure out why the columns for the table are not resizing here
	var smallColWidth int
	if width > 600 {
		smallColWidth = 160
	} else if width > 320 {
		smallColWidth = 60
	} else if width > 240 {
		smallColWidth = 40
	} else if width > 160 {
		smallColWidth = 20
	} else {
		smallColWidth = 14
	}
	columns := []table.Column{
		{Title: "Title 2", Width: width - smallColWidth*2},
		{Title: "Artist", Width: smallColWidth},
		{Title: "Length", Width: smallColWidth},
	}
	tt.table.SetColumns(columns)
}

func NewTestTracksTable() TracksTable {
	columns := []table.Column{
		{Title: "#", Width: 3},
		{Title: "Title", Width: 48},
		{Title: "Artist", Width: 30},
		{Title: "Length", Width: 11},
	}

	rows := []table.Row{
		{"", "Loading...", "", ""},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(8),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)
	t.SetStyles(s)

	tt := TracksTable{table: t, focused: true}

	return tt
}

// func NewTracksTable() TracksTable {
// 	l := list.New([]list.Item{}, itemDelegate{}, defaultWidth, listHeight)
// 	l.Title = "What do you want for dinner?"
// 	l.SetShowStatusBar(false)
// 	l.SetFilteringEnabled(false)
// 	l.SetShowHelp(false)
// 	l.Styles.Title = titleStyle
// 	l.Styles.PaginationStyle = paginationStyle
// 	l.Styles.HelpStyle = helpStyle
//
// 	tt := TracksTable{list: l, focused: true}
//
// 	return tt
// }

package common

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	tableActiveContainerStyle   = BorderContainer()
	tableInactiveContainerStyle = BorderContainer().BorderForeground(lipgloss.Color("243"))
)

var DefaultTracksTableColumns = []table.Column{
	{Title: "#", Width: 3},
	{Title: "Title", Width: 48},
	{Title: "Artist", Width: 30},
	{Title: "Length", Width: 11},
}

var DefaultTracksTableRows = []table.Row{
	{"", "", "", ""},
}

type TracksTable struct {
	focused   bool
	isLoading bool
	onSelect  func(track Track) tea.Cmd
	table     table.Model
	tracks    []Track
}

func NewTracksTable(
	tableHeight int,
	onSelect func(track Track) tea.Cmd,
) TracksTable {
	t := table.New(
		table.WithColumns(DefaultTracksTableColumns),
		table.WithRows(DefaultTracksTableRows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.Bold(false).Foreground(lipgloss.Color("#FFFFFF"))
	t.SetStyles(s)

	tt := TracksTable{table: t, isLoading: false, focused: false, onSelect: onSelect}

	return tt
}

// -------------------------
// --- Lifecycle Methods ---
// -------------------------
func (tt TracksTable) Init() tea.Cmd {
	return nil
}

func (tt TracksTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case " ":
			return tt, nil
		case "enter":
			track := tt.tracks[tt.table.Cursor()]
			cmd := tt.onSelect(track)
			cmds = append(cmds, cmd, PlayTracksCmd(tt.tracks, tt.table.Cursor()))
		}
	}

	tt.table, cmd = tt.table.Update(msg)
	cmds = append(cmds, cmd)

	return tt, tea.Batch(cmds...)
}

func (tt TracksTable) View() string {
	if tt.Focused() {
		return tableActiveContainerStyle.Render(tt.table.View())
	} else {
		return tableInactiveContainerStyle.Render(tt.table.View())
	}
}

// ----------------------------
// --- Tracks Table Methods ---
// ----------------------------
func (tt *TracksTable) SetCursor(idx int) {
	tt.table.SetCursor(idx)
}

func (tt *TracksTable) Focus() {
	tt.focused = true
	tt.table.Focus()

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
	tt.table.SetStyles(s)
}

func (tt *TracksTable) Blur() {
	tt.focused = false
	tt.table.Blur()

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.Bold(false).Foreground(lipgloss.Color("#FFFFFF"))
	tt.table.SetStyles(s)
}

func (tt *TracksTable) Focused() bool {
	return tt.focused
}

func (tt *TracksTable) SetIsLoading(val bool) {
	tt.isLoading = val
	if val {
		var loadingRows = []table.Row{
			{"", "loading...", "", ""},
		}

		tt.table.SetRows(loadingRows)
	}
}

func (tt *TracksTable) IsLoading() bool {
	return tt.isLoading
}

func (tt *TracksTable) UpdateTracks(tracks []Track) {
	var trackRows []table.Row

	for i, track := range tracks {
		trackRows = append(
			trackRows,
			table.Row{
				fmt.Sprint(i + 1),
				RemoveEmojis(track.Title),
				RemoveEmojis(track.User.Name),
				GetLengthText(track.Duration),
				// fmt.Sprint(track.Play_count),
			},
		)
	}

	tt.tracks = tracks
	tt.table.SetRows(trackRows)
}

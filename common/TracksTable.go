package common

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	tableContainerStyle         = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	tableActiveContainerStyle   = tableContainerStyle.Copy().BorderForeground(lipgloss.Color("62"))
	tableInactiveContainerStyle = tableContainerStyle.Copy().BorderForeground(lipgloss.Color("243"))
)

var DefaultTracksTableColumns = []table.Column{
	{Title: "#", Width: 3},
	{Title: "Title", Width: 48},
	{Title: "Artist", Width: 30},
	{Title: "Length", Width: 11},
}

type TracksTable struct {
	focused   bool
	isLoading bool
	onSelect  func(track Track) tea.Cmd
	table     table.Model
	tracks    []Track
}

func NewTracksTable(
	columns []table.Column,
	rows []table.Row,
	onSelect func(track Track) tea.Cmd,
) TracksTable {
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

	tt := TracksTable{table: t, focused: true, onSelect: onSelect}

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
		switch keypress := msg.String(); keypress {
		case "enter":
			track := tt.tracks[tt.table.Cursor()]
			cmd := tt.onSelect(track)
			cmds = append(cmds, cmd, PlayTrackCmd(track))
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
func (tt *TracksTable) Focus() {
	tt.focused = true
}

func (tt *TracksTable) Blur() {
	tt.focused = false
}

func (tt *TracksTable) Focused() bool {
	return tt.focused
}

func (tt *TracksTable) SetIsLoading(val bool) {
	tt.isLoading = val
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
			},
		)
	}

	tt.tracks = tracks
	tt.table.SetRows(trackRows)
}

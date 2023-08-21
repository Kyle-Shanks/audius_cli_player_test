package search

import (
	"app1/api"
	"app1/common"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchResultsMsg struct {
	tracks []common.Track
}

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

func fetchSearchResultsCmd(query string) tea.Cmd {
	return func() tea.Msg {
		tracks, err := api.GetSearchTracks(query)

		if err != nil {
			return errMsg{err}
		}

		return searchResultsMsg{
			tracks: tracks,
		}
	}
}

/* Enum to track search view focus */
type searchFocus = int

const (
	inputFocus searchFocus = iota
	tableFocus
)

type SearchView struct {
	focused bool
	focus   searchFocus
	input   textinput.Model
	table   common.TracksTable
}

func NewSearchView() SearchView {
	i := textinput.New()
	i.Placeholder = "Search"
	i.Focus()

	tt := common.NewTracksTable(
		10,
		func(track common.Track) tea.Cmd {
			return nil
		},
	)

	return SearchView{
		focused: false,
		focus:   inputFocus,
		input:   i,
		table:   tt,
	}
}

func (s SearchView) Init() tea.Cmd {
	return textinput.Blink
}

func (s SearchView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.focused {
			switch key := msg.String(); key {
			case "enter":
				if s.focus == inputFocus {
					val := s.input.Value()
					cmds = append(cmds, fetchSearchResultsCmd(val))
					s.table.SetCursor(0)
					s.table.SetIsLoading(true)
				}
			case "tab":
				if s.focus == inputFocus {
					s.FocusTable()
					return s, nil
				} else {
					s.FocusInput()
					return s, textinput.Blink
				}
			}
		}

	case searchResultsMsg:
		tracks := msg.tracks
		s.table.UpdateTracks(tracks)
		s.FocusTable()
		s.table.SetIsLoading(false)
	}

	if s.focused {
		if s.focus == inputFocus {
			s.input, cmd = s.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			res, cmd := s.table.Update(msg)
			s.table = res.(common.TracksTable)
			cmds = append(cmds, cmd)
		}
	}

	return s, tea.Batch(cmds...)
}

func (s SearchView) View() string {
	inputBorderColor := "62"
	if s.focus == tableFocus {
		inputBorderColor = "242"
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Width(100).Align(lipgloss.Center).Render(
			common.Header().Render("Search"),
		),
		common.BorderContainer().
			BorderForeground(lipgloss.Color(inputBorderColor)).
			Width(100).
			Padding(0, 2).
			Render(
				s.input.View(),
			),
		s.table.View(),
	)
}

func (s *SearchView) Focus() {
	s.focused = true
}

func (s *SearchView) Blur() {
	s.focused = false
}

func (s *SearchView) Focused() bool {
	return s.focused
}

func (s *SearchView) FocusInput() {
	s.focus = inputFocus
	s.input.Focus()
	s.table.Blur()
}

func (s SearchView) InputFocused() bool {
	return s.focus == inputFocus
}

func (s *SearchView) FocusTable() {
	s.focus = tableFocus
	s.input.Blur()
	s.table.Focus()
}

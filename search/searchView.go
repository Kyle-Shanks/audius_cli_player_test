package search

import tea "github.com/charmbracelet/bubbletea"

type SearchView struct {
}

func NewSearchView() SearchView {
	return SearchView{}
}

func (s SearchView) Init() tea.Cmd {
	return nil
}

func (s SearchView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) {}

	var cmd tea.Cmd

	return s, cmd
}

func (s SearchView) View() string {
	return "Search View"
}

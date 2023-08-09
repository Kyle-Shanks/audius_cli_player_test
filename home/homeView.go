package home

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeView struct {
	sideList       SimpleList
	trendingTracks TrendingTracks
}

func NewHomeView() HomeView {
	return HomeView{
		sideList:       NewTestSimpleList(),
		trendingTracks: NewTrendingTracks(),
	}
}

func (h HomeView) Init() tea.Cmd {
	return tea.Batch(
		h.sideList.Init(),
		h.trendingTracks.Init(),
	)
}

func (h HomeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// var res tea.Model
	// res, cmd = h.sideList.Update(msg)
	// h.sideList = res.(SimpleList)

	var res tea.Model
	res, cmd = h.trendingTracks.Update(msg)
	h.trendingTracks = res.(TrendingTracks)

	return h, cmd
}

func (h HomeView) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		h.trendingTracks.View(),
	)
}

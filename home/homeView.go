package home

import (
	// "fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* Enum to track which view to render */
type homeFocus = int

const (
	sideList homeFocus = iota
	mainContent
	bottomContent
)

type HomeView struct {
	focus    homeFocus
	sideList SimpleList
	table    TracksTable
}

func NewHomeView() HomeView {
	return HomeView{
		sideList: NewTestSimpleList(),
		table:    NewTestTracksTable(),
		focus:    sideList,
	}
}

func (h HomeView) Init() tea.Cmd {
	// Initialize sub-models
	return tea.Batch(
		h.sideList.Init(),
		h.table.Init(),
	)
}

func (h HomeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Add any home view msg handling here
	// switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	fmt.Println(msg.Width, msg.Height)
	// }

	var cmd tea.Cmd

	switch h.focus {
	case sideList:
		var res tea.Model
		res, cmd = h.sideList.Update(msg)
		h.sideList = res.(SimpleList)
	}

	var res tea.Model
	res, cmd = h.table.Update(msg)
	h.table = res.(TracksTable)

	return h, cmd
}

func (h HomeView) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		h.table.View(),
	)
}

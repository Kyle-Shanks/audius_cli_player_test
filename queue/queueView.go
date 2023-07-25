package queue

import tea "github.com/charmbracelet/bubbletea"

type QueueView struct {
}

func NewQueueView() QueueView {
	return QueueView{}
}

func (q QueueView) Init() tea.Cmd {
	return nil
}

func (q QueueView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) {}

	var cmd tea.Cmd

	return q, cmd
}

func (q QueueView) View() string {
	return "Queue View"
}

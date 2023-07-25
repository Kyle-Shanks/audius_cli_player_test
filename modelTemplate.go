package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type TestModel struct {
	// Add stateand submodels here
}

func (tm TestModel) Init() tea.Cmd {
	return nil
}

func (tm TestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return tm, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "":
			return tm, tea.Quit
		}
	}

	// Pass msg to appropriate submodels here

	return tm, cmd
}

func (tm TestModel) View() string {
	// Render out the needed submodels and things here
	return ""
}

func NewTestModel() TestModel {
	tm := TestModel{}
	// Add submodels and initial state

	return tm
}

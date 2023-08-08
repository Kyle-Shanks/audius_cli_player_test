package main

import (
	"fmt"
	"os"

	"app1/home"
	"app1/player"
	"app1/queue"
	"app1/search"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* Enum to track which view to render */
type appView = int

const (
	homeView appView = iota
	searchView
	queueView
)

// TODO: Add better doc comments
type App struct {
	/* Current view that the app is displaying */
	view appView
	/* Now playing view */
	player player.Player
	/* Home view of the app */
	homeView home.HomeView
	/* Queue view of the app */
	queueView queue.QueueView
	/* Search view of the app */
	searchView search.SearchView

	/* App key map */
	keyMap KeyMap
	/* App help text component */
	help help.Model
}

/* Create and initialize a new instance of the App */
func NewApp() App {

	h := help.New()
	h.Styles.ShortDesc = h.Styles.ShortDesc.Foreground(lipgloss.Color("#555"))
	h.Styles.FullDesc = h.Styles.FullDesc.Foreground(lipgloss.Color("#555"))
	h.Styles.ShortKey = h.Styles.ShortKey.Foreground(lipgloss.Color("#777"))
	h.Styles.FullKey = h.Styles.FullKey.Foreground(lipgloss.Color("#777"))

	app := App{
		view:       homeView,
		player:     player.NewPlayer(),
		homeView:   home.NewHomeView(),
		queueView:  queue.NewQueueView(),
		searchView: search.NewSearchView(),

		keyMap: AppKeyMap,
		help:   h,
	}

	return app
}

func (a App) Init() tea.Cmd {
	// Initialize sub-models
	return tea.Batch(
		a.player.Init(),
		a.homeView.Init(),
		a.queueView.Init(),
		a.searchView.Init(),
	)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.resizeApp(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.keyMap.Help):
			a.help.ShowAll = !a.help.ShowAll
		case key.Matches(msg, a.keyMap.Quit):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, a.keyMap.Play):
			a.player.TogglePause()
		case key.Matches(msg, a.keyMap.Mute):
			a.player.ToggleMute()
		}
	}

	var updateRes tea.Model

	// Pass msg to active view
	switch a.view {
	case homeView:
		updateRes, cmd = a.homeView.Update(msg)
		a.homeView = updateRes.(home.HomeView)
		cmds = append(cmds, cmd)
	case searchView:
		updateRes, cmd = a.searchView.Update(msg)
		a.searchView = updateRes.(search.SearchView)
		cmds = append(cmds, cmd)
	case queueView:
		updateRes, cmd = a.queueView.Update(msg)
		a.queueView = updateRes.(queue.QueueView)
		cmds = append(cmds, cmd)
	}

	// Pass msg to Player
	updateRes, cmd = a.player.Update(msg)
	a.player = updateRes.(player.Player)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	var mainView tea.Model

	switch a.view {
	case homeView:
		mainView = a.homeView
	case queueView:
		mainView = a.queueView
	case searchView:
		mainView = a.searchView
	default:
		mainView = a.homeView
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		mainView.View(),
		a.player.View(),
		a.getHelpText(),
	)
}

func (a App) getHelpText() string {
	helpContainerStyle := lipgloss.NewStyle().Width(100).Align(lipgloss.Left)

	return helpContainerStyle.Render(
		a.help.View(a.keyMap),
	)
}

// Run when a window size message is received
func (a *App) resizeApp(width int, height int) {
	// fmt.Println(width, height)
}

func main() {
	program := tea.NewProgram(NewApp(), tea.WithAltScreen())
	// TODO: Add tea.WithMouseAllMotion() later to add mouse click handling

	if _, err := program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

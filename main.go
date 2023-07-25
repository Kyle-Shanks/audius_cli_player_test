package main

import (
	"fmt"
	"os"

	"app1/home"
	"app1/player"
	"app1/queue"
	"app1/search"

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
	player *player.Player
	/* Home view of the app */
	homeView home.HomeView
	/* Queue view of the app */
	queueView queue.QueueView
	/* Search view of the app */
	searchView search.SearchView
	/* Test error */
	err error
}

/* Create and initialize a new instance of the App */
func NewApp() App {
	app := App{
		player:     player.NewPlayer(),
		homeView:   home.NewHomeView(),
		queueView:  queue.NewQueueView(),
		searchView: search.NewSearchView(),
	}

	return app
}

type statusMsg string

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

func testHTTP() tea.Msg {
	track, err := GetTrackById("JpvZ0")

	if err != nil {
		return errMsg{err}
	}

	// entries, err := os.ReadDir(".")
	// var fileNames []string
	//
	// for _, entry := range entries {
	// 	fileNames = append(fileNames, entry.Name())
	// }
	//
	// fmt.Println(entries)
	// fmt.Println(fileNames)

	n, err := GetTrackMp3("JpvZ0")

	fmt.Println(n)

	fmt.Printf("%q by %v", track.Title, track.User.Name)

	return statusMsg(track.Title)
}

func (a App) Init() tea.Cmd {
	// return testHTTP
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Process msg
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.resizeApp(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmds = append(cmds, tea.Quit)
		case "a":
			a.player.Play("./audio/QmS6k7iWF3BnmdaQC5taJb2yhPy5grKtCvsQoWwEfM6nQp.mp3")
			// cmds = append(cmds, player.TrackPlayMsg)
			// a.audioPlayer.Play(os.TempDir()+"tempTrack.mp31834572435")
		case "p":
			a.player.TogglePause()
		case "m":
			a.player.ToggleMute()
		}
	case statusMsg:
		// fmt.Println(msg)
	case errMsg:
		a.err = msg
		// fmt.Println(msg)
	}

	var updateRes tea.Model

	// Pass msg to active view
	switch a.view {
	case homeView:
		updateRes, cmd = a.homeView.Update(msg)
		a.homeView = updateRes.(home.HomeView)
		cmds = append(cmds, cmd)
	case queueView:
		updateRes, cmd = a.queueView.Update(msg)
		a.queueView = updateRes.(queue.QueueView)
		cmds = append(cmds, cmd)
	case searchView:
		updateRes, cmd = a.searchView.Update(msg)
		a.searchView = updateRes.(search.SearchView)
		cmds = append(cmds, cmd)
	}

	// Pass msg to now playing
	// updateRes, cmd = a.player.Update(msg)
	// a.player = updateRes.(*player.Player)
	// cmds = append(cmds, cmd)

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
		lipgloss.Top,
		mainView.View(),
		a.player.View(),
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

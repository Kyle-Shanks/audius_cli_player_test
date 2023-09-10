package main

import (
	"fmt"
	"os"

	"app1/api"
	"app1/common"
	"app1/player"
	"app1/queue"
	"app1/search"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* Enum to track which view to render */
type appView = int

const (
	trendingView appView = iota
	undergroundView
	favoritesView
	searchView
	queueView
)

// TODO: Add better doc comments
type App struct {
	/* Current view that the app is displaying */
	view appView
	/* Now playing view */
	player player.Player
	/* Trending view of the app */
	trendingView common.TracksTableView
	/* Underground view of the app */
	undergroundView common.TracksTableView
	/* Favorites view of the app */
	favoritesView common.TracksTableView
	/* Queue view of the app */
	queueView queue.QueueView
	/* Search view of the app */
	searchView search.SearchView

	/* App key map */
	keyMap KeyMap
	/* App help text component */
	help help.Model
}

func getMyFavs() ([]common.Track, error) {
	return api.GetUserFavoriteTracks("PWVbN")
}

/* Create and initialize a new instance of the App */
func NewApp() App {
	h := help.New()
	h.Styles.ShortDesc = h.Styles.ShortDesc.Foreground(lipgloss.Color("#555"))
	h.Styles.FullDesc = h.Styles.FullDesc.Foreground(lipgloss.Color("#555"))
	h.Styles.ShortKey = h.Styles.ShortKey.Foreground(lipgloss.Color("#777"))
	h.Styles.FullKey = h.Styles.FullKey.Foreground(lipgloss.Color("#777"))

	app := App{
		view:   trendingView,
		player: player.NewPlayer(),
		trendingView: common.NewTracksTableView(
			"Trending Tracks",
			api.GetTrendingTracks,
		),
		undergroundView: common.NewTracksTableView(
			"Underground Tracks",
			api.GetUndergroundTracks,
		),
		favoritesView: common.NewTracksTableView(
			"Favorites",
			getMyFavs,
		),
		queueView:  queue.NewQueueView(),
		searchView: search.NewSearchView(),

		keyMap: AppKeyMap,
		help:   h,
	}
	app.trendingView.Focus()

	return app
}

func (a App) Init() tea.Cmd {
	// Initialize sub-models
	return tea.Batch(
		a.player.Init(),
		a.trendingView.Init(),
		a.undergroundView.Init(),
		a.favoritesView.Init(),
		a.queueView.Init(),
		a.searchView.Init(),
	)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	searchInputFocused := a.view == searchView && a.searchView.InputFocused()

	// TODO: Add Queue view key bind handler
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.resizeApp(msg.Width, msg.Height)

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			// log := fmt.Sprintf("x:%v, y:%v", msg.X, msg.Y)
			// cmds = append(cmds, common.LogCmd(log))

			// Switch view/tab
			if msg.Y == 1 {
				if msg.X >= 11 && msg.X <= 31 {
					a.updateViewFocus(trendingView)
				} else if msg.X >= 34 && msg.X <= 57 {
					a.updateViewFocus(undergroundView)
				} else if msg.X >= 60 && msg.X <= 74 {
					a.updateViewFocus(favoritesView)
				} else if msg.X >= 77 && msg.X <= 88 {
					a.updateViewFocus(searchView)
					cmds = append(cmds, textinput.Blink)
				}
			}
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.keyMap.Help):
			if !searchInputFocused {
				a.help.ShowAll = !a.help.ShowAll
			}
		case key.Matches(msg, a.keyMap.Quit):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, a.keyMap.Underground):
			if !searchInputFocused {
				a.updateViewFocus(undergroundView)
				return a, tea.Batch(cmds...)
			}
		case key.Matches(msg, a.keyMap.Trending):
			if !searchInputFocused {
				a.updateViewFocus(trendingView)
				return a, tea.Batch(cmds...)
			}
		case key.Matches(msg, a.keyMap.Favorites):
			if !searchInputFocused {
				a.updateViewFocus(favoritesView)
				return a, tea.Batch(cmds...)
			}
		case key.Matches(msg, a.keyMap.Search):
			if !searchInputFocused {
				a.updateViewFocus(searchView)
				if msg.String() == "/" {
					a.searchView.FocusInput()
				}
				cmds = append(cmds, textinput.Blink)
				return a, tea.Batch(cmds...)
			}
		}

		// Exit early if in search input
		// TODO: Find a better way to handle this
		if searchInputFocused {
			updateRes, cmd := a.searchView.Update(msg)
			a.searchView = updateRes.(search.SearchView)
			cmds = append(cmds, cmd)

			return a, tea.Batch(cmds...)
		}
	}

	var updateRes tea.Model

	// Pass msg to views
	updateRes, cmd = a.trendingView.Update(msg)
	a.trendingView = updateRes.(common.TracksTableView)
	cmds = append(cmds, cmd)

	updateRes, cmd = a.undergroundView.Update(msg)
	a.undergroundView = updateRes.(common.TracksTableView)
	cmds = append(cmds, cmd)

	updateRes, cmd = a.favoritesView.Update(msg)
	a.favoritesView = updateRes.(common.TracksTableView)
	cmds = append(cmds, cmd)

	updateRes, cmd = a.searchView.Update(msg)
	a.searchView = updateRes.(search.SearchView)
	cmds = append(cmds, cmd)

	updateRes, cmd = a.queueView.Update(msg)
	a.queueView = updateRes.(queue.QueueView)
	cmds = append(cmds, cmd)

	// Pass msg to Player
	updateRes, cmd = a.player.Update(msg)
	a.player = updateRes.(player.Player)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	var mainView tea.Model

	trendingHeader := common.Header().MarginLeft(1)
	undergroundHeader := common.Header().MarginLeft(2)
	favoritesHeader := common.Header().MarginLeft(2)
	// queueHeader := common.Header().MarginLeft(2)
	searchHeader := common.Header().MarginLeft(2)

	switch a.view {
	case trendingView:
		mainView = a.trendingView
		trendingHeader.Foreground(lipgloss.Color("229")).Background(common.PrimaryAlt)
	case undergroundView:
		mainView = a.undergroundView
		undergroundHeader.Foreground(lipgloss.Color("229")).Background(common.PrimaryAlt)
	case favoritesView:
		mainView = a.favoritesView
		favoritesHeader.Foreground(lipgloss.Color("229")).Background(common.PrimaryAlt)
	case searchView:
		mainView = a.searchView
		searchHeader.Foreground(lipgloss.Color("229")).Background(common.PrimaryAlt)
	// case queueView:
	// 	mainView = a.queueView
	// 	queueHeader.Foreground(lipgloss.Color("229")).Background(common.PrimaryAlt)
	default:
		mainView = a.trendingView
	}

	headerTabs := lipgloss.NewStyle().
		Align(lipgloss.Center).
		MarginTop(1).
		Width(100).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				trendingHeader.Render("(T)rending Tracks"),
				undergroundHeader.Render("(U)nderground Tracks"),
				favoritesHeader.Render("(F)avorites"),
				// queueHeader.Render("(Q)ueue"),
				searchHeader.Render("(S)earch"),
			),
		)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerTabs,
		mainView.View(),
		a.player.View(),
		a.getHelpText(),
	)
}

func (a App) getHelpText() string {
	helpContainerStyle := lipgloss.NewStyle().Width(100).Align(lipgloss.Left).MarginLeft(2)
	// return helpContainerStyle.Render(
	// 	a.help.View(a.keyMap),
	// )

	if a.help.ShowAll {
		return helpContainerStyle.Render(
			a.help.FullHelpView(append(a.player.KeyMap.FullHelp(), a.keyMap.FullHelp()...)),
		)
	} else {
		return helpContainerStyle.Render(
			a.help.ShortHelpView(a.keyMap.ShortHelp()),
		)
	}
}

func (a *App) updateViewFocus(newView appView) {
	if a.view != newView {
		a.view = newView

		a.trendingView.Blur()
		a.undergroundView.Blur()
		a.favoritesView.Blur()
		a.searchView.Blur()

		switch a.view {
		case trendingView:
			a.trendingView.Focus()
		case undergroundView:
			a.undergroundView.Focus()
		case favoritesView:
			a.favoritesView.Focus()
		case searchView:
			a.searchView.Focus()
		}
	}
}

// Run when a window size message is received
func (a *App) resizeApp(width int, height int) {
	// fmt.Println(width, height)
}

func main() {
	program := tea.NewProgram(
		NewApp(),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)
	// TODO: Add tea.WithMouseAllMotion() later to add mouse click handling

	if _, err := program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

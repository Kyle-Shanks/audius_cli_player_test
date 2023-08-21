package common

import "github.com/charmbracelet/lipgloss"

// TODO: Figure out how to make this a map or something
// -------------------------
// --- Color Definitions ---
// -------------------------
// type Colors struct {
//   white lipgloss.Color
// }

// var Colors = map[string]lipgloss.Color{
// 	"white": lipgloss.Color("#FFFFFF"),
// }

var Grey0 = lipgloss.Color("#1A1A1A")
var Grey1 = lipgloss.Color("#333333")
var Grey2 = lipgloss.Color("#555555")
var Grey3 = lipgloss.Color("#777777")
var Grey4 = lipgloss.Color("#999999")
var Grey5 = lipgloss.Color("#AAAAAA")
var Grey6 = lipgloss.Color("#C2C2C2")
var White = lipgloss.Color("#FFFFFF")

var Primary = lipgloss.Color("62")
var PrimaryAlt = lipgloss.Color("57")
var Inactive = lipgloss.Color("243")
var EmptyColor = lipgloss.Color("#3B4252")

// ---------------------------
// --- Reusable Components ---
// ---------------------------
var header = lipgloss.
	NewStyle().
	Foreground(White).
	Background(EmptyColor).
	Bold(true).
	Padding(0, 2)

func Header() lipgloss.Style {
	return header.Copy()
}

var activeHeader = header.Copy().
	Foreground(lipgloss.Color("229")).
	Background(PrimaryAlt)

func ActiveHeader() lipgloss.Style {
	return activeHeader.Copy()
}

// TODO: Make an active and inactive one
var borderContainer = lipgloss.
	NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Primary)

func BorderContainer() lipgloss.Style {
	return borderContainer.Copy()
}

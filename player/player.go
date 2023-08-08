package player

import (
	"app1/api"
	"app1/common"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	containerStyle         = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")).Padding(0, 4)
	inactiveContainerStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("243")).Padding(0, 4)
)

type TrackPlayMsg struct {
	trackName  string
	artistName string
}

type Player struct {
	audioPlayer  AudioPlayer
	currentTrack common.Track
	currentPos   float64
	progress     progress.Model
}

func NewPlayer() Player {
	prog := progress.New(
		progress.WithDefaultGradient(),
		// progress.WithSolidFill("62"),
		progress.WithoutPercentage(),
		progress.WithWidth(92),
	)
	prog.EmptyColor = "#3B4252"

	// Half size blocks
	prog.Empty = '▄'
	prog.Full = '▄'

	// Full size blocks
	// prog.Empty = '█'
	// prog.Full = '█'

	p := Player{
		audioPlayer:  NewAudioPlayer(),
		currentTrack: common.Track{Title: ""},
		currentPos:   0.0,
		progress:     prog,
	}

	return p
}

type tickMsg time.Time

func (p Player) Init() tea.Cmd {
	return tea.Batch(
		p.progress.Init(),
		tickCmd(),
	)
}

func (p Player) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case common.PlayTrackMsg:
		p.PlayTrack(msg.Track)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			p.audioPlayer.DeleteTempFiles()
		}

	case tickMsg:
		p.UpdateProgressPos()
		cmds = append(cmds, tickCmd())
	}

	// Pass msg to appropriate submodels here
	// p.progress, cmd = p.progress.Update(msg)
	// cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p Player) View() string {
	var text string
	if p.currentTrack.Title != "" {
		text = lipgloss.JoinHorizontal(
			lipgloss.Center,
			lipgloss.NewStyle().Bold(true).Render(
				p.currentTrack.Title,
			),
			lipgloss.NewStyle().Foreground(lipgloss.Color("#929292")).Render(
				" - "+p.currentTrack.User.Name,
			),
		)
	} else {
		text = "Now Playing"
	}

	return containerStyle.Width(100).AlignHorizontal(lipgloss.Center).Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			text,
			p.progress.ViewAs(p.currentPos),
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				lipgloss.NewStyle().Width(46).Align(lipgloss.Left).Render(
					common.GetDurationText(int(p.currentPos*float64(p.currentTrack.Duration))),
				),
				lipgloss.NewStyle().Width(46).Align(lipgloss.Right).Render(
					common.GetDurationText(p.currentTrack.Duration),
				),
			),
		),
	)
}

func (p *Player) UpdateProgressPos() {
	if p.audioPlayer.Streamer != nil {
		pos := p.audioPlayer.Streamer.Position()
		p.currentPos = float64(pos/48000) / float64(p.currentTrack.Duration)
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// - Audio Player Methods -
func (p *Player) PlayTrack(track common.Track) {
	fileName, err := api.GetTrackMp3(track.Id)

	if err != nil {
		fmt.Println(err)
		return
	}

	p.currentTrack = track

	p.Play(fileName)
}

func (p *Player) Play(filepath string) {
	var err error
	_, err = p.audioPlayer.Play(filepath)

	if err != nil {
		fmt.Println(err)
	}

	// trackLen := p.audioPlayer.Play(filepath)
	// return TrackPlayMsg{trackName: }
}

func (p *Player) Pause() {
	p.audioPlayer.Pause()
}

func (p *Player) Resume() {
	p.audioPlayer.Resume()
}

func (p *Player) TogglePause() {
	p.audioPlayer.TogglePause()
}

func (p *Player) SetVolume(vol float64) {
	p.audioPlayer.SetVolume(vol)
}

func (p *Player) Mute() {
	p.audioPlayer.Mute()
}

func (p *Player) Unmute() {
	p.audioPlayer.Unmute()
}

func (p *Player) ToggleMute() {
	p.audioPlayer.ToggleMute()
}

func (p *Player) Seek(pos int) {
	p.audioPlayer.Seek(pos)
}

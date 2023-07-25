package player

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	containerStyle         = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	inactiveContainerStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("243"))
)

type TrackPlayMsg struct {
	trackName  string
	artistName string
}

type Player struct {
	// TODO: Need to get track type in here
	// currentTrack Track

	// Add state and submodels here
	audioPlayer *AudioPlayer
}

func NewPlayer() *Player {
	p := &Player{
		audioPlayer: NewAudioPlayer(),
	}

	return p
}

func (p Player) Init() tea.Cmd {
	return nil
}

func (p Player) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return p, nil

	case tea.KeyMsg:
		switch msg.String() {
		// case "":
		// 	return p, tea.Quit
		}
	}

	// Pass msg to appropriate submodels here

	return p, cmd
}

func (p Player) View() string {
	// Render out the needed submodels and things here
	return containerStyle.Render("Now Playing")
}

// - Audio Player Methods -
func (p *Player) Play(filepath string) {
	p.audioPlayer.Play(filepath)
	// trackLen := p.audioPlayer.Play(filepath)
	// return TrackPlayMsg{trackName: }
}

func (p *Player) Pause() {
	if p.audioPlayer.Ctrl != nil {
		p.audioPlayer.Pause()
	}
}

func (p *Player) Resume() {
	if p.audioPlayer.Ctrl != nil {
		p.audioPlayer.Resume()
	}
}

func (p *Player) TogglePause() {
	if p.audioPlayer.Ctrl != nil {
		p.audioPlayer.TogglePause()
	}
}

func (p *Player) SetVolume(vol float64) {
	if p.audioPlayer.Volume != nil {
		p.audioPlayer.SetVolume(vol)
	}
}

func (p *Player) Mute() {
	if p.audioPlayer.Volume != nil {
		p.audioPlayer.Mute()
	}
}

func (p *Player) Unmute() {
	if p.audioPlayer.Volume != nil {
		p.audioPlayer.Unmute()
	}
}

func (p *Player) ToggleMute() {
	if p.audioPlayer.Volume != nil {
		p.audioPlayer.ToggleMute()
	}
}

func (p *Player) Seek(pos int) {
	if p.audioPlayer.Streamer != nil {
		p.audioPlayer.Seek(pos)
	}
}

package player

import (
	"fmt"
	"github.com/Kyle-Shanks/audius_cli_player_test/api"
	"github.com/Kyle-Shanks/audius_cli_player_test/common"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	containerStyle         = common.BorderContainer().Padding(0, 2)
	inactiveContainerStyle = containerStyle.Copy().BorderForeground(common.Inactive)
)

type FetchTrackMp3ResMsg struct {
	track    common.Track
	fileName string
}

type Player struct {
	audioPlayer  AudioPlayer
	tracksQueue  []common.Track
	queuePos     int
	currentTrack common.Track
	currentPos   float64
	muted        bool
	paused       bool
	repeat       bool
	progress     progress.Model
	KeyMap       KeyMap
}

func NewPlayer() Player {
	prog := progress.New(
		progress.WithDefaultGradient(),
		// progress.WithSolidFill("62"),
		progress.WithoutPercentage(),
		progress.WithWidth(96),
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
		tracksQueue:  []common.Track{},
		queuePos:     0,
		currentTrack: common.Track{Title: ""},
		currentPos:   0.0,
		muted:        false,
		paused:       false,
		repeat:       false,
		progress:     prog,
		KeyMap:       PlayerKeyMap,
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
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			// Handle progress bar seeking
			if p.audioPlayer.Streamer != nil && msg.Y == 21 && (msg.X >= 3 && msg.X <= 98) {
				percentage := float32(msg.X-3) / 96.0
				pos := int(float32(p.audioPlayer.Streamer.Len()) * percentage)
				p.Seek(pos)
			}
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, p.KeyMap.Pause):
			p.TogglePause()
		case key.Matches(msg, p.KeyMap.Mute):
			p.ToggleMute()
		case key.Matches(msg, p.KeyMap.Repeat):
			p.ToggleRepeat()
		// case key.Matches(msg, p.KeyMap.VolumeUp):
		// 	fmt.Println("1")
		// case key.Matches(msg, p.KeyMap.VolumeDown):
		// 	fmt.Println("2")
		// case key.Matches(msg, p.KeyMap.SkipForward):
		// 	pos := p.audioPlayer.Streamer.Position()
		// 	p.Seek(pos + (10 * 48000))
		// case key.Matches(msg, p.KeyMap.SkipBack):
		// 	pos := p.audioPlayer.Streamer.Position()
		// 	if pos <= 10*48000 {
		// 		p.Seek(0)
		// 	} else {
		// 		p.Seek(pos - (10 * 48000))
		// 	}
		case key.Matches(msg, p.KeyMap.Quit):
			p.audioPlayer.DeleteTempFiles()
		}

	case common.PlayTrackMsg:
		p.tracksQueue = []common.Track{msg.Track}
		p.queuePos = 0
		cmds = append(cmds, FetchAndPlayTrackCmd(p.tracksQueue[p.queuePos]))

	case common.PlayTracksMsg:
		p.tracksQueue = msg.Tracks
		p.queuePos = msg.QueuePos
		cmds = append(cmds, FetchAndPlayTrackCmd(p.tracksQueue[p.queuePos]))

	case FetchTrackMp3ResMsg:
		p.currentTrack = msg.track
		p.play(msg.fileName)

	case tickMsg:
		p.UpdateProgressPos()
		if TrackEnded {
			cmds = append(cmds, p.handlePlaybackEnd())
		}
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
				common.RemoveEmojis(p.currentTrack.Title),
			),
			lipgloss.NewStyle().Foreground(common.Grey4).Render(
				" - "+common.RemoveEmojis(p.currentTrack.User.Name),
			),
		)
	} else {
		text = common.Header().Render("Now Playing")
	}

	var repeatText string
	var pauseText string
	var muteText string

	if p.repeat {
		repeatText = lipgloss.NewStyle().Foreground(lipgloss.Color("#77F")).Render("repeat")
	} else {
		repeatText = lipgloss.NewStyle().Foreground(common.Grey3).Render("repeat")
	}

	if p.paused {
		pauseText = lipgloss.NewStyle().Foreground(lipgloss.Color("#7F7")).Render("pause")
	} else {
		pauseText = lipgloss.NewStyle().Foreground(common.Grey3).Render("pause")
	}

	if p.muted {
		muteText = lipgloss.NewStyle().Foreground(lipgloss.Color("#F77")).Render("mute")
	} else {
		muteText = lipgloss.NewStyle().Foreground(common.Grey3).Render("mute")
	}

	return containerStyle.Width(100).AlignHorizontal(lipgloss.Center).Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			text,
			p.progress.ViewAs(p.currentPos),
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				lipgloss.NewStyle().Width(32).Align(lipgloss.Left).Render(
					common.GetDurationText(int(p.currentPos*float64(p.currentTrack.Duration))),
				),
				lipgloss.NewStyle().Width(32).Align(lipgloss.Center).Foreground(common.White).Render(
					fmt.Sprintf("%v  %v  %v", repeatText, pauseText, muteText),
				),
				lipgloss.NewStyle().Width(32).Align(lipgloss.Right).Render(
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

func (p *Player) handlePlaybackEnd() tea.Cmd {
	TrackEnded = false

	if p.repeat {
		p.play(p.audioPlayer.CurrentTrackFileName)
		return nil
	} else {
		p.audioPlayer.DeleteTempFiles()

		if p.queuePos+1 < len(p.tracksQueue) {
			p.queuePos++
			return FetchAndPlayTrackCmd(p.tracksQueue[p.queuePos])
		}

		return nil
	}
}

// - Player Commands -
/* Command to fetch track mp3 and return saved mp3 filename */
func FetchAndPlayTrackCmd(track common.Track) tea.Cmd {
	return func() tea.Msg {
		fileName, err := api.GetTrackMp3(track.Id)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		return FetchTrackMp3ResMsg{track: track, fileName: fileName}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// - Audio Player Methods -
func (p *Player) SetAudioPlayer(ap AudioPlayer) {
	p.audioPlayer = ap
}

func (p *Player) play(filepath string) {
	p.paused = false
	_, err := p.audioPlayer.Play(filepath, p.muted)

	if err != nil {
		fmt.Println(err)
	}
}

func (p *Player) Pause() {
	if p.audioPlayer.Ctrl != nil {
		p.audioPlayer.Pause()
		p.paused = p.audioPlayer.Ctrl.Paused
	}
}

func (p *Player) Resume() {
	if p.audioPlayer.Ctrl != nil {
		p.audioPlayer.Resume()
		p.paused = p.audioPlayer.Ctrl.Paused
	}
}

func (p *Player) TogglePause() {
	if p.audioPlayer.Ctrl != nil {
		p.audioPlayer.TogglePause()
		p.paused = p.audioPlayer.Ctrl.Paused
	}
}

func (p *Player) SetVolume(vol float64) {
	p.audioPlayer.SetVolume(vol)
}

func (p *Player) Mute() {
	if p.audioPlayer.Volume != nil {
		p.audioPlayer.Mute()
		p.muted = p.audioPlayer.Volume.Silent
	} else {
		p.muted = true
	}
}

func (p *Player) Unmute() {
	if p.audioPlayer.Volume != nil {
		p.audioPlayer.Unmute()
		p.muted = p.audioPlayer.Volume.Silent
	} else {
		p.muted = false
	}
}

func (p *Player) ToggleMute() {
	if p.audioPlayer.Volume != nil {
		p.audioPlayer.ToggleMute()
		p.muted = p.audioPlayer.Volume.Silent
	} else {
		p.muted = !p.muted
	}
}

func (p *Player) ToggleRepeat() {
	p.repeat = !p.repeat
}

func (p *Player) Seek(pos int) {
	p.audioPlayer.Seek(pos)
}

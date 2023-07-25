package player

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

const DefaultSampleRate = 44100

func DecodeFile(filename string, file *os.File) (beep.StreamSeekCloser, beep.Format, error) {
	switch {
	case strings.HasSuffix(filename, ".mp3"):
		return mp3.Decode(file)
	case strings.HasSuffix(filename, ".wav"):
		return wav.Decode(file)
	case strings.HasSuffix(filename, ".flac"):
		return flac.Decode(file)
	case strings.HasSuffix(filename, ".ogg"):
		return vorbis.Decode(file)
	default:
		return mp3.Decode(file)
	}
}

type AudioPlayer struct {
	Ctrl       *beep.Ctrl
	Streamer   beep.StreamSeekCloser
	Volume     *effects.Volume
	SampleRate int
}

func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{
		SampleRate: DefaultSampleRate,
	}
}

func OnPlaybackEnd() {
	fmt.Println("Playback ended")
	// TODO: Handle looping and queuing the next song and stuff here
}

func (ap *AudioPlayer) Play(filepath string) int {
	fmt.Println("Playing " + filepath)
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}

	streamer, format, err := DecodeFile(file.Name(), file)
	if err != nil {
		fmt.Println(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// NOTE: beep.Loop tested here, but want to loop manually to allow user to toggle easily
	// ap.Ctrl = &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
	ap.Ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	ap.Streamer = streamer
	ap.Volume = &effects.Volume{
		Streamer: ap.Ctrl,
		Base:     2,
		Volume:   -3,
		Silent:   false,
	}
	ap.SampleRate = int(format.SampleRate)

	speaker.Play(beep.Seq(ap.Volume, beep.Callback(OnPlaybackEnd)))

	// Return the track length
	return ap.Streamer.Len() / int(format.SampleRate)
}

func (ap AudioPlayer) GetTrackLength() int {
	return ap.Streamer.Len() / ap.SampleRate
}

// TODO: Add error handling to all of these methods
func (ap AudioPlayer) Pause() {
	speaker.Lock()
	ap.Ctrl.Paused = true
	speaker.Unlock()
}

func (ap AudioPlayer) Resume() {
	speaker.Lock()
	ap.Ctrl.Paused = false
	speaker.Unlock()
}

func (ap AudioPlayer) TogglePause() {
	if ap.Ctrl.Paused {
		ap.Resume()
	} else {
		ap.Pause()
	}
}

func (ap AudioPlayer) Mute() {
	ap.Volume.Silent = true
}

func (ap AudioPlayer) Unmute() {
	ap.Volume.Silent = false
}

func (ap AudioPlayer) ToggleMute() {
	if ap.Volume.Silent {
		ap.Unmute()
	} else {
		ap.Mute()
	}
}

func (ap AudioPlayer) SetVolume(vol float64) {
	ap.Volume.Volume = vol
}

func (ap AudioPlayer) Seek(pos int) {
	speaker.Lock()
	ap.Streamer.Seek(pos)
	speaker.Unlock()
}

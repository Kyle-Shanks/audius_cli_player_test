package player

import (
	// "fmt"
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

// const DefaultSampleRate = 44100
const DefaultSampleRate = 48000

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
	Ctrl                 *beep.Ctrl
	Streamer             beep.StreamSeekCloser
	Volume               *effects.Volume
	sampleRate           int
	currentTrackFileName string
}

func NewAudioPlayer() AudioPlayer {
	return AudioPlayer{
		sampleRate:           DefaultSampleRate,
		currentTrackFileName: "",
	}
}

func (ap AudioPlayer) OnPlaybackEnd() {
	// fmt.Println("Playback ended")
	ap.DeleteTempFiles()

	// TODO: Handle looping and queueing the next song and stuff here
}

func (ap *AudioPlayer) DeleteTempFiles() {
	if ap.currentTrackFileName != "" {
		os.Remove(ap.currentTrackFileName)
	}
}

func (ap *AudioPlayer) Play(filepath string) (*AudioPlayer, error) {
	// fmt.Println("Playing " + filepath)
	// Delete prev temp files
	ap.DeleteTempFiles()

	ap.currentTrackFileName = filepath
	file, err := os.Open(filepath)
	if err != nil {
		return ap, err
	}

	streamer, format, err := DecodeFile(file.Name(), file)
	if err != nil {
		return ap, err
	}

	// TODO: Should make a separate init function for the speaker
	// Can just use default as the sampleRate
	if ap.Streamer == nil {
		// Initialize
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	} else {
		// Already played music, reset
		ap.Streamer.Close()
		speaker.Clear()
	}

	// NOTE: beep.Loop tested here, but want to loop manually like in next line to allow user to toggle easily
	// ap.Ctrl = &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
	ap.Ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	ap.Streamer = streamer
	ap.Volume = &effects.Volume{
		Streamer: ap.Ctrl,
		Base:     2,
		Volume:   -3,
		Silent:   false,
	}
	ap.sampleRate = int(format.SampleRate)
	// fmt.Println(ap.sampleRate)

	speaker.Play(beep.Seq(ap.Volume, beep.Callback(ap.OnPlaybackEnd)))

	// Return the track length
	// return ap.Streamer.Len() / int(format.SampleRate)
	return ap, nil
}

func (ap AudioPlayer) GetTrackLength() int {
	return ap.Streamer.Len() / ap.sampleRate
}

// TODO: Add error handling to all of these methods
func (ap *AudioPlayer) Pause() {
	if ap.Ctrl != nil {
		speaker.Lock()
		ap.Ctrl.Paused = true
		speaker.Unlock()
	}
}

func (ap *AudioPlayer) Resume() {
	if ap.Ctrl != nil {
		speaker.Lock()
		ap.Ctrl.Paused = false
		speaker.Unlock()
	}
}

func (ap *AudioPlayer) TogglePause() {
	if ap.Ctrl != nil {
		if ap.Ctrl.Paused {
			ap.Resume()
		} else {
			ap.Pause()
		}
	}
}

func (ap *AudioPlayer) Mute() {

	if ap.Volume != nil {
		ap.Volume.Silent = true
	}
}

func (ap *AudioPlayer) Unmute() {
	if ap.Volume != nil {
		ap.Volume.Silent = false
	}
}

func (ap *AudioPlayer) ToggleMute() {
	if ap.Volume != nil {
		if ap.Volume.Silent {
			ap.Unmute()
		} else {
			ap.Mute()
		}
	}
}

func (ap *AudioPlayer) SetVolume(vol float64) {
	if ap.Volume != nil {
		ap.Volume.Volume = vol
	}
}

func (ap *AudioPlayer) Seek(pos int) {
	if ap.Streamer != nil {
		speaker.Lock()
		ap.Streamer.Seek(pos)
		speaker.Unlock()
	}
}

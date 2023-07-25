package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// const AUDIUS_API_ENDPOINT string = "https://api.audius.co/"
const AUDIUS_API_ENDPOINT string = "https://discoveryprovider.audius.co/v1"

type HttpErrMsg struct{ error }

func (e HttpErrMsg) Error() string { return e.error.Error() }

// --- Helper Functions ---
/* Create and send request and return the byte array */
func get(path string) ([]byte, error) {
	// Create Request
	req, _ := http.NewRequest("GET", AUDIUS_API_ENDPOINT+path, nil)

	// Add query params
	q := req.URL.Query()
	q.Add("app_name", "audius-cli")
	req.URL.RawQuery = q.Encode()

	fmt.Println("GET Request Url: " + req.URL.String())

	// Add headers
	req.Header.Set("Accept", "application/json")

	// Create client and submit request
	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()
	return io.ReadAll(res.Body)
}

/* Send get request and unmarshal track data */
func getTrack(trackPath string) (Track, error) {
	resBytes, err := get(trackPath)

	if err != nil {
		return Track{}, err
	}

	var trackRes TrackResponse
	err = json.Unmarshal([]byte(resBytes), &trackRes)

	return trackRes.Data, err
}

/* Send get request and unmarshal array of tracks */
func getTracks(tracksPath string) ([]Track, error) {
	resBytes, err := get(tracksPath)

	if err != nil {
		return []Track{}, err
	}

	var tracksRes TracksResponse
	err = json.Unmarshal([]byte(resBytes), &tracksRes)

	return tracksRes.Data, err
}

// --- Get Functions ---
func GetTrackById(trackId string) (Track, error) {
	path := "/tracks/" + trackId
	return getTrack(path)
}

func GetUserTracks(userId string) ([]Track, error) {
	path := "/users/" + userId + "/tracks"
	return getTracks(path)
}

func GetPlaylistTracks(playlistId string) ([]Track, error) {
	path := "/playlists/" + playlistId + "/tracks"
	return getTracks(path)
}

func GetTrendingTracks() ([]Track, error) {
	path := "/tracks/trending"
	return getTracks(path)
}

/* Fetch track mp3 and return file name */
func GetTrackMp3(trackId string) (string, error) {
	path := "/tracks/" + trackId + "/stream"
	resBytes, err := get(path)

	if err != nil {
		return "", err
	}

	// TODO: Update to delete the file later
	// Also figure out if you can label by trackId so you can reuse files
	file, err := os.CreateTemp(os.TempDir(), "tempTrack.*.mp3")
	fmt.Println(os.TempDir())
	file.Write(resBytes)

	return file.Name(), err
}

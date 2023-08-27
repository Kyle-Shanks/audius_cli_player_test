package api

import (
	"app1/common"
	"encoding/json"
	// "fmt"
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

	// fmt.Println("GET Request Url: " + req.URL.String())

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
func getTrack(trackPath string) (common.Track, error) {
	resBytes, err := get(trackPath)

	if err != nil {
		return common.Track{}, err
	}

	var trackRes common.TrackApiResponse
	err = json.Unmarshal([]byte(resBytes), &trackRes)

	return trackRes.Data, err
}

/* Send get request and unmarshal array of tracks */
func getTracks(tracksPath string) ([]common.Track, error) {
	resBytes, err := get(tracksPath)

	if err != nil {
		return []common.Track{}, err
	}

	var tracksRes common.TracksApiResponse
	err = json.Unmarshal([]byte(resBytes), &tracksRes)

	return tracksRes.Data, err
}

// --- Get Functions ---
// Track Functions
func GetTrackById(trackId string) (common.Track, error) {
	path := "/tracks/" + trackId
	return getTrack(path)
}

func GetUserTracks(userId string) ([]common.Track, error) {
	path := "/users/" + userId + "/tracks"
	return getTracks(path)
}

func GetPlaylistTracks(playlistId string) ([]common.Track, error) {
	path := "/playlists/" + playlistId + "/tracks"
	return getTracks(path)
}

func GetTrendingTracks() ([]common.Track, error) {
	path := "/tracks/trending"
	return getTracks(path)
}

func GetUndergroundTracks() ([]common.Track, error) {
	path := "/tracks/trending/underground"
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
	file, err := os.CreateTemp(os.TempDir(), "TEMP_APP_TRACK.*.mp3")
	// fmt.Println(os.TempDir())
	file.Write(resBytes)

	return file.Name(), err
}

func GetSearchTracks(query string) ([]common.Track, error) {
	// Create Request
	path := "/tracks/search"
	req, _ := http.NewRequest("GET", AUDIUS_API_ENDPOINT+path, nil)

	// Add query params
	q := req.URL.Query()
	q.Add("app_name", "audius-cli")
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	// Add headers
	req.Header.Set("Accept", "application/json")

	// Create client and submit request
	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return []common.Track{}, err
	}

	defer res.Body.Close()
	resBytes, err := io.ReadAll(res.Body)

	if err != nil {
		return []common.Track{}, err
	}

	var tracksRes common.TracksApiResponse
	err = json.Unmarshal([]byte(resBytes), &tracksRes)

	return tracksRes.Data, err
}

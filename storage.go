package main

import (
	"encoding/json"
	"fmt"
	"github.com/senseyeio/duration"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type yt_video struct {
	UploadDate  time.Time `json:"upload_date"`
	Duration    int64     `json:"duration"`
	Title       string    `json:"title"`
	ID          string    `json:"id"`
	Description string    `json:"description"`
}

type yt_playlist struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

const YT_DATA_LOC = "data/youtube/"
const YT_STOR_LOC = "data/youtube_storage/"
const VIDEO_ABR = 128

func download_yt_video(id string) string {
	fmt.Printf("Get request for youtube video %s\n", id)
	os.MkdirAll(YT_STOR_LOC, 0777)

	video_path := YT_STOR_LOC + id + ".m4a"
	if _, err := os.Stat(video_path); os.IsNotExist(err) {
		fmt.Printf("Downloading video %s to %s\n", id, video_path)

		out, _ := exec.Command("youtube-dl",
			"-x",
			"--audio-format",
			"mp3",
			"https://www.youtube.com/watch?v="+id,
			"-o",
			video_path).Output()
		fmt.Println(string(out))
	}

	return video_path
}

func get_playlist_data(id string) *yt_playlist {
	var playlist_data yt_playlist

	fmt.Printf("Get metadata for playlist %s\n", id)

	// Create folder if it doesn't exist
	os.MkdirAll(YT_DATA_LOC, 0777)

	// Check if playlist metadata exists
	playlist_path := YT_DATA_LOC + "/playlist-" + id + ".json"
	if _, err := os.Stat(playlist_path); os.IsNotExist(err) {
		fmt.Printf("Playlist metadata doesn't exist\n")

		// Ask the API layer for playlist data
		response := api_get_playlist_data(id)

		playlist_data.Name = response.Items[0].Snippet.Title
		playlist_data.Description = response.Items[0].Snippet.Description

		// Save it do disk
		file, _ := os.Create(playlist_path)
		playlist_json, _ := json.MarshalIndent(playlist_data, "", " ")
		file.Write(playlist_json)
	} else {
		data, _ := ioutil.ReadFile(playlist_path)
		_ = json.Unmarshal(data, &playlist_data)
	}

	return &playlist_data
}

func get_yt_video_data(id string) *yt_video {
	fmt.Printf("Get metadata for video ID %s\n", id)
	var video_data yt_video

	// Create folder if it doesn't exist
	os.MkdirAll(YT_DATA_LOC, 0777)

	// Check if video metadata exists
	video_path := YT_DATA_LOC + "/video-" + id + ".json"
	if _, err := os.Stat(video_path); os.IsNotExist(err) {
		fmt.Printf("Metadata doesn't exist\n")

		// Ask the API layer for video data
		response := api_get_video_data(id)

		// Store it in a structure
		video_data.UploadDate, _ = time.Parse(time.RFC3339, response.Items[0].Snippet.PublishedAt)
		// Parse the duration
		dur, _ := duration.ParseISO8601(response.Items[0].ContentDetails.Duration)
		video_data.Duration = int64(dur.TS + dur.TM*60 + dur.TH*60*60 + dur.D*60*60*24)
		video_data.Title = response.Items[0].Snippet.Title
		video_data.ID = response.Items[0].Id
		video_data.Description = response.Items[0].Snippet.Description

		// Save it to disk
		file, _ := os.Create(video_path)
		video_json, _ := json.MarshalIndent(video_data, "", " ")
		file.Write(video_json)

	} else {
		data, _ := ioutil.ReadFile(video_path)
		_ = json.Unmarshal(data, &video_data)
	}

	return &video_data
}

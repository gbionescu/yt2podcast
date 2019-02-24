package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type configdata struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
}

var cfg_data configdata

// Get XML for a channel
func get_podcast(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	xml := api_get_yt_channel(vars["yt_channel"])
	fmt.Fprintf(w, string(xml))
}

// Video request entry point
func get_video(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	video_path := download_yt_video(vars["video_id"])
	fmt.Printf("Serving %s\n", video_path)
	http.ServeFile(w, r, video_path)
}

// Get XML for a playlist
func get_playlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fmt.Printf("Get playlist %s\n", vars["playlist_id"])
	xml := api_get_yt_playlist(vars["playlist_id"])

	fmt.Fprintf(w, string(xml))
}

// Load config from disk
func load_cfg(path string) {
	data, _ := ioutil.ReadFile(path)
	_ = json.Unmarshal(data, &cfg_data)
}

// Returns the port on which the server is running on
func get_port() string {
	return cfg_data.Port
}

// Get the address where the podcast is running
func get_podcast_addr() string {
	return cfg_data.Hostname + ":" + cfg_data.Port
}

func main() {
	load_cfg("config.json")

	r := mux.NewRouter()
	r.HandleFunc("/api/ytchan/{yt_channel}", get_podcast)
	r.HandleFunc("/api/ytplaylist/{playlist_id}", get_playlist)
	r.HandleFunc("/api/ytv/{video_id}", get_video)

	server := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:" + string(get_port()),
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  300 * time.Second,
	}
	fmt.Printf("Listening on %s\n", get_port())
	log.Fatal(server.ListenAndServe())
}

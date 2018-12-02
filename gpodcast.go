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

func get_podcast(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	request_type := vars["type"]
	if request_type != "user" && request_type != "channel" {
		return
	}

	xml := get_yt_podcast(request_type + "/" + vars["yt_channel"])
	fmt.Fprintf(w, string(xml))
}

func get_video(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	request_type := vars["type"]
	if request_type != "user" && request_type != "channel" {
		return
	}

	fmt.Printf("video get channel: %s, ID: %s\n", vars["yt_channel"], vars["video_id"])

	video_path := get_yt_video(request_type+"/"+vars["yt_channel"], vars["video_id"])
	fmt.Printf("Serving %s\n", video_path)
	http.ServeFile(w, r, video_path)
}

func load_cfg(path string) {
	data, err := ioutil.ReadFile(path)
	check_panic(err)
	err = json.Unmarshal(data, &cfg_data)
	check_panic(err)
}

func get_port() string {
	return cfg_data.Port
}

func get_podcast_addr() string {
	return cfg_data.Hostname + ":" + cfg_data.Port
}

func main() {
	load_cfg("config.json")

	r := mux.NewRouter()
	r.HandleFunc("/podcast/youtube/{type}/{yt_channel}", get_podcast)
	r.HandleFunc("/podcast/youtube-video/{type}/{yt_channel}/{video_id}", get_video)

	server := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:" + string(get_port()),
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  300 * time.Second,
	}
	fmt.Printf("Listening on %s\n", get_port())
	log.Fatal(server.ListenAndServe())
}

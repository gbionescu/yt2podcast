package main

import (
    "log"
    "net/http"
    "fmt"
    "time"
    "encoding/json"
    "io/ioutil"
    "github.com/gorilla/mux"
)

type configdata struct {
    Hostname string `json:"hostname"`
    Port     string `json:"port"`
}

var cfg_data configdata

func get_podcast(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    
    xml := get_yt_podcast(vars["yt_channel"])
    fmt.Fprintf(w, string(xml))
}

func get_video(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    fmt.Printf("video get %s %s", vars["yt_channel"], vars["video_id"]);
    
    //video_id := get_yt_video(vars["yt_channel"], vars["video_id"])
    //fmt.Fprintf(w, xml)
}

func load_cfg(path string) {
    data, err := ioutil.ReadFile(path)
    check_panic(err)
    err = json.Unmarshal(data, &cfg_data)
    check_panic(err)
}

func get_port() (string) {
    return cfg_data.Port
}

func get_podcast_addr() (string) {
    return cfg_data.Hostname + ":" + cfg_data.Port
}

func main() {
    load_cfg("config.json")
    
    r := mux.NewRouter()
    r.HandleFunc("/podcast/youtube/{yt_channel}", get_podcast)
    r.HandleFunc("/podcast/youtube-video/{yt_channel}/{video_id}", get_video)

    server := &http.Server{
        Handler:      r,
        Addr:         "0.0.0.0:" + string(get_port()),
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    fmt.Printf("Listening on %s\n", get_port())
    log.Fatal(server.ListenAndServe())
}

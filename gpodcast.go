package main

import (
    "flag"
    "log"
    "net/http"
    "fmt"
    "html"
    "time"
    "github.com/gorilla/mux"
    "os/exec"
)

func get_podcast(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    
    // Get youtube channel URL
    var yt_chan_url = fmt.Sprintf("https://www.youtube.com/channel/%s", vars["yt_channel"])
    
    // Get JSON containing all channel information
    fmt.Printf("Running youtube-dl for %s\n", yt_chan_url)
    out, err := exec.Command("youtube-dl", "-J", yt_chan_url).Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s\n", out)
    
    fmt.Fprintf(w, "%q", html.EscapeString(out))
}

func main() {
    port := flag.String("p", "8100", "port to serve on")
    flag.Parse()
    
    r := mux.NewRouter()
    r.HandleFunc("/podcast/youtube/{yt_channel}", get_podcast)

    server := &http.Server{
        Handler:      r,
        Addr:         "0.0.0.0:" + *port,
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    log.Fatal(server.ListenAndServe())
}
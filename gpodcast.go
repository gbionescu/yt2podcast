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
    "encoding/json"
    "os"
)

const YT_CHAN_LOC="data/youtube/"

type yt_video struct {
    Upload_date string  `json:"upload_date"`
    Title string        `json:"title"`
    Id string           `json:"id"`
    Duration int        `json:"duration"`
    Description string  `json:"description"`
}

type yt_channel struct {
    Uploader string `json:"uploader"`
    Entries []yt_video `json:"entries"`
}

type yt_metadata struct {
    Chan_data yt_channel
    Last_request time.Time
}

func get_podcast(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    // Get youtube channel URL
    var yt_chan_url = fmt.Sprintf("https://www.youtube.com/channel/%s", vars["yt_channel"])

    // Check if channel exists and load the specific data
    yt_chan_path := YT_CHAN_LOC + vars["yt_channel"]
    yt_chan_data := yt_chan_path + "/chan-data.json"
    os.MkdirAll(yt_chan_path, 0777)

    var chan_data yt_metadata
    chan_data.Last_request = time.Now()

    if _, err := os.Stat(yt_chan_data); os.IsNotExist(err) {
        fmt.Printf("XXX\n");
    }

    // Get JSON containing all channel information
    fmt.Printf("Running youtube-dl for %s\n", yt_chan_url)
    out, err := exec.Command("youtube-dl", "-J", yt_chan_url).Output()
    if err != nil {
        log.Fatal(err)
    }

    // Decode the JSON
    var json_data yt_channel
    if err := json.Unmarshal(out, &json_data); err != nil {
        panic(err)
    }

    fmt.Printf("uploader %s\n", json_data.Uploader);

    for _, entry := range json_data.Entries {
        //fmt.Printf("%#v\n", val);
        fmt.Printf("title %s\n", entry.Title);
    }

    fmt.Fprintf(w, "%q", html.EscapeString(string(out)))

    if f, err := os.Create(yt_chan_data); err != nil {
        panic(err);
    }

    f.Write(json.Marshal(chan_data));

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
    fmt.Printf("Listening on %s\n", *port)
    log.Fatal(server.ListenAndServe())
}

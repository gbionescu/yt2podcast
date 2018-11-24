package main

import (
    "flag"
    "log"
    "net/http"
    "fmt"
    "time"
    "github.com/gorilla/mux"
    "os/exec"
    "encoding/json"
    "os"
    "io/ioutil"
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

func check_panic(e error) {
    if e != nil {
        panic(e)
    }
}

func get_podcast(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    // Get youtube channel URL
    var yt_chan_url = fmt.Sprintf("https://www.youtube.com/channel/%s", vars["yt_channel"])

    // Check if channel exists and load the specific data
    yt_chan_path := YT_CHAN_LOC + vars["yt_channel"]
    yt_chan_data := yt_chan_path + "/chan-data.json"
    os.MkdirAll(yt_chan_path, 0777)

    var metadata yt_metadata

    if _, err := os.Stat(yt_chan_data); os.IsNotExist(err) {
        fmt.Printf("Channel does not exist. Will create it.\n");
    } else {
        fmt.Printf("Reading data from existing folder.\n");
        data, err := ioutil.ReadFile(yt_chan_data)
        check_panic(err)

        err = json.Unmarshal(data, &metadata)
        check_panic(err)
    }
    
    last_request := metadata.Last_request;

    // Get JSON containing all channel information
    fmt.Printf("Running youtube-dl for %s\n", yt_chan_url)
    
    if last_request != 0 {
        fmt.Printf("Last request was: %s\n", last_request)
    }
    
    out, err := exec.Command("youtube-dl", "-J", yt_chan_url).Output()
    check_panic(err)

    // Decode the JSON
    var json_data yt_channel
    err = json.Unmarshal(out, &json_data)
    check_panic(err)

    // Fill the metadata structure
    metadata.Chan_data.Uploader = json_data.Uploader;
    metadata.Last_request = time.Now()

    for _, entry := range json_data.Entries {
        fmt.Printf("title %s\n", entry.Title);
    }

    fmt.Fprintf(w, "Done\n")

    f, err := os.Create(yt_chan_data)
    check_panic(err)

    chan_json, err := json.MarshalIndent(metadata, "", "  ")
    check_panic(err)
    
    f.Write(chan_json)
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

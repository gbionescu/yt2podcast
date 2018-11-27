package main


import (
    "fmt"
    "time"
    "os"
    "os/exec"
    "encoding/json"
    "io/ioutil"
    "bytes"
    "github.com/jbub/podcasts"
    "strconv"
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

func get_yt_podcast(yt_chan string) (string) {

    // Get youtube channel URL
    var yt_chan_url = fmt.Sprintf("https://www.youtube.com/channel/%s", yt_chan)

    // Check if channel exists and load the specific data
    yt_chan_path := YT_CHAN_LOC + yt_chan
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
    
    if !last_request.IsZero() {
        fmt.Printf("Last request was: %s\n", last_request)
        
        dateafter := fmt.Sprintf("%d%d%d", 
            last_request.Year(), last_request.Month(), last_request.Day())
        
        // Only get videos uploaded after the last request
        fmt.Printf("Adding to youtube-dl: %s\n", dateafter)
        out, err := exec.Command("youtube-dl", "-j", "--dateafter", dateafter, yt_chan_url).Output()
        check_panic(err)
        
        // Per-video info is split by '\n'
        videos := bytes.Split(out, []byte("\n"));
        
        // Parse each new video
        for video_no := range videos {
            var entry yt_video

            if len(videos[video_no]) == 0 {
                continue
            }
            
            err := json.Unmarshal(videos[video_no], &entry)
            check_panic(err)
            
            fmt.Printf("Checking title %s\n", entry.Title);
            
            // Check for duplicates
            exists := false
            for _, crt_entry := range metadata.Chan_data.Entries {
                if entry.Id == crt_entry.Id {
                    exists = true
                }
            }
            
            if !exists {
                fmt.Printf("-> Adding to channel.\n");
                metadata.Chan_data.Entries = append(metadata.Chan_data.Entries, entry);
            }
        }
    } else {
        out, err := exec.Command("youtube-dl", "-J", yt_chan_url).Output()
        check_panic(err)
        
        // Decode the JSON
        var fetched_data yt_channel
        err = json.Unmarshal(out, &fetched_data)
        check_panic(err)

        // Fill the metadata structure
        metadata.Chan_data.Uploader = fetched_data.Uploader;
        metadata.Chan_data.Entries = fetched_data.Entries
    }

    metadata.Last_request = time.Now()

    f, err := os.Create(yt_chan_data)
    check_panic(err)

    chan_json, err := json.MarshalIndent(metadata, "", "  ")
    check_panic(err)
    
    f.Write(chan_json)
    
    p := &podcasts.Podcast{
        Title:       metadata.Chan_data.Uploader,
        Description: "",
        Language:    "EN",
    }
    
    for _, crt_entry := range metadata.Chan_data.Entries {
        
        upload_date, err := time.Parse("20060102", crt_entry.Upload_date)
        check_panic(err)
        
        duration := strconv.FormatInt(int64(crt_entry.Duration), 10)
        
        p.AddItem(&podcasts.Item{
            Title:   crt_entry.Title,
            PubDate:  &podcasts.PubDate{upload_date},
            Enclosure: &podcasts.Enclosure{
                URL:    "http://www.example-podcast.com/my-podcast/2/episode.mp3",
                Length: duration,
                Type:   "MP3",
            },
        })
    }
    
    feed, err := p.Feed()
    check_panic(err)
    
    xml, err := feed.XML()
    
    return xml
}

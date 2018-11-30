package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eduncan911/podcast"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type yt_video struct {
	PlaylistUploader string      `json:"playlist_uploader"`
	UploadDate       string      `json:"upload_date"`
	Extractor        string      `json:"extractor"`
	Series           interface{} `json:"series"`
	Format           string      `json:"format"`
	Vbr              interface{} `json:"vbr"`
	Chapters         interface{} `json:"chapters"`
	Height           int         `json:"height"`
	LikeCount        int         `json:"like_count"`
	Duration         int         `json:"duration"`
	Fulltitle        string      `json:"fulltitle"`
	PlaylistIndex    int         `json:"playlist_index"`
	RequestedFormats []struct {
		HTTPHeaders struct {
			AcceptCharset  string `json:"Accept-Charset"`
			AcceptLanguage string `json:"Accept-Language"`
			AcceptEncoding string `json:"Accept-Encoding"`
			Accept         string `json:"Accept"`
			UserAgent      string `json:"User-Agent"`
		} `json:"http_headers"`
		Tbr               float64 `json:"tbr"`
		Protocol          string  `json:"protocol"`
		Format            string  `json:"format"`
		URL               string  `json:"url"`
		Vcodec            string  `json:"vcodec"`
		FormatNote        string  `json:"format_note"`
		Height            int     `json:"height,omitempty"`
		DownloaderOptions struct {
			HTTPChunkSize int `json:"http_chunk_size"`
		} `json:"downloader_options"`
		Width     int    `json:"width,omitempty"`
		Ext       string `json:"ext"`
		Filesize  int    `json:"filesize"`
		Fps       int    `json:"fps,omitempty"`
		FormatID  string `json:"format_id"`
		PlayerURL string `json:"player_url"`
		Quality   int    `json:"quality"`
		Acodec    string `json:"acodec"`
		Abr       int    `json:"abr,omitempty"`
	} `json:"requested_formats"`
	ViewCount          int         `json:"view_count"`
	Playlist           string      `json:"playlist"`
	Title              string      `json:"title"`
	Filename           string      `json:"_filename"`
	Creator            interface{} `json:"creator"`
	Ext                string      `json:"ext"`
	ID                 string      `json:"id"`
	DislikeCount       int         `json:"dislike_count"`
	PlaylistID         string      `json:"playlist_id"`
	Abr                int         `json:"abr"`
	UploaderURL        string      `json:"uploader_url"`
	Categories         []string    `json:"categories"`
	Fps                int         `json:"fps"`
	StretchedRatio     interface{} `json:"stretched_ratio"`
	SeasonNumber       interface{} `json:"season_number"`
	Annotations        interface{} `json:"annotations"`
	WebpageURLBasename string      `json:"webpage_url_basename"`
	Acodec             string      `json:"acodec"`
	DisplayID          string      `json:"display_id"`
	AutomaticCaptions  struct {
	} `json:"automatic_captions"`
	Description        string        `json:"description"`
	Tags               []interface{} `json:"tags"`
	Track              interface{}   `json:"track"`
	RequestedSubtitles interface{}   `json:"requested_subtitles"`
	StartTime          interface{}   `json:"start_time"`
	AverageRating      float64       `json:"average_rating"`
	Uploader           string        `json:"uploader"`
	FormatID           string        `json:"format_id"`
	EpisodeNumber      interface{}   `json:"episode_number"`
	UploaderID         string        `json:"uploader_id"`
	Subtitles          struct {
	} `json:"subtitles"`
	PlaylistTitle string `json:"playlist_title"`
	Thumbnails    []struct {
		URL string `json:"url"`
		ID  string `json:"id"`
	} `json:"thumbnails"`
	License      interface{} `json:"license"`
	Artist       interface{} `json:"artist"`
	ExtractorKey string      `json:"extractor_key"`
	Vcodec       string      `json:"vcodec"`
	AltTitle     interface{} `json:"alt_title"`
	Thumbnail    string      `json:"thumbnail"`
	ChannelID    string      `json:"channel_id"`
	IsLive       interface{} `json:"is_live"`
	EndTime      interface{} `json:"end_time"`
	WebpageURL   string      `json:"webpage_url"`
	Formats      []struct {
		HTTPHeaders struct {
			AcceptCharset  string `json:"Accept-Charset"`
			AcceptLanguage string `json:"Accept-Language"`
			AcceptEncoding string `json:"Accept-Encoding"`
			Accept         string `json:"Accept"`
			UserAgent      string `json:"User-Agent"`
		} `json:"http_headers"`
		FormatNote        string  `json:"format_note"`
		Protocol          string  `json:"protocol"`
		Format            string  `json:"format"`
		URL               string  `json:"url"`
		Vcodec            string  `json:"vcodec"`
		Tbr               float64 `json:"tbr,omitempty"`
		Abr               int     `json:"abr,omitempty"`
		PlayerURL         string  `json:"player_url"`
		DownloaderOptions struct {
			HTTPChunkSize int `json:"http_chunk_size"`
		} `json:"downloader_options,omitempty"`
		Ext        string `json:"ext"`
		Filesize   int    `json:"filesize,omitempty"`
		FormatID   string `json:"format_id"`
		Quality    int    `json:"quality,omitempty"`
		Acodec     string `json:"acodec"`
		Container  string `json:"container,omitempty"`
		Height     int    `json:"height,omitempty"`
		Width      int    `json:"width,omitempty"`
		Fps        int    `json:"fps,omitempty"`
		Resolution string `json:"resolution,omitempty"`
	} `json:"formats"`
	PlaylistUploaderID string      `json:"playlist_uploader_id"`
	ChannelURL         string      `json:"channel_url"`
	Resolution         interface{} `json:"resolution"`
	Width              int         `json:"width"`
	NEntries           int         `json:"n_entries"`
	AgeLimit           int         `json:"age_limit"`
}
type yt_channel struct {
	Extractor          string     `json:"extractor"`
	Type               string     `json:"_type"`
	Uploader           string     `json:"uploader"`
	Entries            []yt_video `json:"entries"`
	ID                 string     `json:"id"`
	Title              string     `json:"title"`
	ExtractorKey       string     `json:"extractor_key"`
	UploaderID         string     `json:"uploader_id"`
	UploaderURL        string     `json:"uploader_url"`
	WebpageURL         string     `json:"webpage_url"`
	WebpageURLBasename string     `json:"webpage_url_basename"`
}

type yt_metadata struct {
	Chan_data    yt_channel
	Last_request time.Time
}

const YT_CHAN_LOC = "data/youtube/"

func check_panic(e error) {
	if e != nil {
		panic(e)
	}
}

func get_yt_podcast(yt_chan string) []byte {

	// Get youtube channel URL
	var yt_chan_url = fmt.Sprintf("https://www.youtube.com/channel/%s", yt_chan)

	// Check if channel exists and load the specific data
	yt_chan_path := YT_CHAN_LOC + yt_chan
	yt_chan_data := yt_chan_path + "/chan-data.json"
	os.MkdirAll(yt_chan_path, 0777)

	var metadata yt_metadata

	if _, err := os.Stat(yt_chan_data); os.IsNotExist(err) {
		fmt.Printf("Channel does not exist. Will create it.\n")
	} else {
		fmt.Printf("Reading data from existing folder.\n")
		data, err := ioutil.ReadFile(yt_chan_data)
		check_panic(err)

		err = json.Unmarshal(data, &metadata)
		check_panic(err)
	}

	last_request := metadata.Last_request

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
		videos := bytes.Split(out, []byte("\n"))

		// Parse each new video
		for video_no := range videos {
			var entry yt_video

			if len(videos[video_no]) == 0 {
				continue
			}

			err := json.Unmarshal(videos[video_no], &entry)
			check_panic(err)

			fmt.Printf("Checking title %s\n", entry.Title)

			// Check for duplicates
			exists := false
			for _, crt_entry := range metadata.Chan_data.Entries {
				if entry.ID == crt_entry.ID {
					exists = true
				}
			}

			if !exists {
				fmt.Printf("-> Adding to channel.\n")
				metadata.Chan_data.Entries = append(metadata.Chan_data.Entries, entry)
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
		metadata.Chan_data.Uploader = fetched_data.Uploader
		metadata.Chan_data.Entries = fetched_data.Entries
	}

	metadata.Last_request = time.Now()

	f, err := os.Create(yt_chan_data)
	check_panic(err)

	chan_json, err := json.MarshalIndent(metadata, "", "  ")
	check_panic(err)

	f.Write(chan_json)

	p := podcast.New(
		metadata.Chan_data.Uploader,
		"title",
		"description",
		nil,
		nil)

	p.IExplicit = "no"

	fmt.Printf("Generating XML for %d entries", len(metadata.Chan_data.Entries))

	for _, crt_entry := range metadata.Chan_data.Entries {
		upload_date, err := time.Parse("20060102", crt_entry.UploadDate)
		check_panic(err)

		link := "http://" + get_podcast_addr() + "/podcast/youtube-video/" + crt_entry.ChannelID + "/" + crt_entry.ID
		item := podcast.Item{
			Title:       "Title: " + crt_entry.Title,
			Link:        link,
			Description: "Description: " + crt_entry.Description,
			PubDate:     &upload_date,
		}

		item.AddDuration(int64(crt_entry.Duration))

		fmt.Printf("Checking filesize for %s\n", crt_entry.Title)
		fsize := 0
		for _, format := range crt_entry.Formats {
			if format.Ext == "m4a" && format.Abr == 128 {
				fsize = format.Filesize
				fmt.Printf("Found as: %d\n", fsize)
				break
			}
		}

		item.AddEnclosure(link, podcast.MP3, int64(fsize))

		_, err = p.AddItem(item)
		check_panic(err)
	}

	return p.Bytes()
}

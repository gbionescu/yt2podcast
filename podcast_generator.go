package main

import (
	"fmt"
	"github.com/eduncan911/podcast"
)

// Generate podcast XML using a given playlist ID and a list of youtube videos
func gen_xml_yt_playlist(plist_id string, list []string) string {
	plist_data := get_playlist_data(plist_id)

	fmt.Printf("Generating XML for playlist ID %s\n", plist_id)

	// Create podcast instance
	p := podcast.New(
		plist_data.Name,
		"link",
		plist_data.Description,
		nil,
		nil)
	p.IExplicit = "no"

	for _, crt_entry := range list {
		// Get video metadata
		video_data := get_yt_video_data(crt_entry)

		// Generate the link
		link := "http://" + get_podcast_addr() +
			"/api/ytv/" + video_data.ID

		description := video_data.Description
		if description == "" {
			description = "No description"
		}

		item := podcast.Item{
			Title:       video_data.Title,
			Link:        link,
			Description: description,
			PubDate:     &video_data.UploadDate,
		}

		item.AddDuration(video_data.Duration)
		item.AddEnclosure(link, podcast.MP3, 1024)

		_, _ = p.AddItem(item)
	}

	return string(p.Bytes())
}

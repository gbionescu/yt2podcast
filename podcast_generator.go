package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/eduncan911/podcast"
	"strings"
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

		title := strings.Replace(video_data.Title, "%", "%%", -1)

		description := video_data.Description
		if description == "" {
			description = "No description"
		}
		description = strings.Replace(description, "%", "%%", -1)

		title_escaped := bytes.NewBufferString("")
		descr_escaped := bytes.NewBufferString("")

		xml.EscapeText(title_escaped, []byte(title))
		xml.EscapeText(descr_escaped, []byte(description))

		item := podcast.Item{
			Title:       title_escaped.String(),
			Link:        link,
			Description: descr_escaped.String(),
			PubDate:     &video_data.UploadDate,
		}

		item.AddDuration(video_data.Duration)
		item.AddEnclosure(link, podcast.MP3, 1024)

		_, _ = p.AddItem(item)
	}

	return string(p.Bytes())
}

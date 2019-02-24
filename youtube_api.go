package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

// Get the upload playlist ID for a channel
func get_channel_uploads_id(service *youtube.Service, part string, id string) string {
	fmt.Printf("Checking %s\n", id)

	call := service.Channels.List(part)
	fmt.Printf("Checking if it's a channel ID\n")
	call = call.Id(id)
	response, err := call.Do()
	handleError(err, "")

	if len(response.Items) == 0 {
		call = service.Channels.List(part)
		fmt.Printf("Checking it it's a user\n")
		call = call.ForUsername(id)
		response, err = call.Do()
		handleError(err, "")
	}

	return response.Items[0].ContentDetails.RelatedPlaylists.Uploads
}

func get_service() *youtube.Service {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.New(client)
	handleError(err, "Error creating YouTube client")

	return service
}

// Get XML for a youtube channel - returns XML for uploads
func api_get_yt_channel(channel string) string {
	fmt.Printf("API get youtube channel %s\n", channel)

	// Get the uploads playlist ID
	var service = get_service()
	uploads_id := get_channel_uploads_id(service, "snippet,contentDetails,status", channel)

	return api_get_yt_playlist(uploads_id)
}

// Gets videos listed on a playlist page
// Needed because youtube limits a request to maximum 50 videos
func get_playlist_page(service *youtube.Service, part string, id string, page string) *youtube.PlaylistItemListResponse {
	call := service.PlaylistItems.List(part)

	// Get maximum number of items
	call = call.MaxResults(50)
	call = call.PlaylistId(id)

	// Check if we should request a page
	if page != "" {
		call = call.PageToken(page)
	}

	response, err := call.Do()

	handleError(err, "")

	return response
}

// Generate XML for a playlist ID
func api_get_yt_playlist(id string) string {
	fmt.Printf("API get playlist ID %s\n", id)

	var video_list []string
	service := get_service()

	// Go through each page and collect the videos
	// TODO cache content to make less requests
	next_page_token := ""
	for {
		// Get a playlist page and add video IDs to the list
		response := get_playlist_page(service, "snippet,contentDetails", id, next_page_token)

		for _, item := range response.Items {
			video_list = append(video_list, item.ContentDetails.VideoId)
		}

		// Keep going while there is a next page
		if response.NextPageToken == "" {
			break
		}
		next_page_token = response.NextPageToken
	}

	return gen_xml_yt_playlist(id, video_list)
}

// Gets data for a given playlist ID
func api_get_playlist_data(id string) *youtube.PlaylistListResponse {
	service := get_service()
	call := service.Playlists.List("snippet,contentDetails")

	// Set max results to 0 because we only want playlist metadata
	call = call.MaxResults(0)

	call = call.Id(id)

	response, _ := call.Do()

	return response
}

// Gets data for a given youtube video ID
func api_get_video_data(id string) *youtube.VideoListResponse {
	fmt.Printf("Getting video data for %s\n", id)

	service := get_service()
	call := service.Videos.List("snippet,contentDetails,statistics")
	call = call.Id(id)

	response, err := call.Do()
	handleError(err, "")

	return response
}

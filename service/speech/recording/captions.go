package recording

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func handleError(err error, str string) {
	if err != nil {
		panic(err)
	}
}

// Retrieve playlistItems in the specified playlist
func playlistItemsList(service *youtube.Service, part []string, playlistId string, pageToken string) *youtube.PlaylistItemListResponse {
	call := service.PlaylistItems.List(part)
	call = call.PlaylistId(playlistId)
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
	response, err := call.Do()
	handleError(err, "")
	return response
}

// Retrieve resource for the authenticated user's channel
func channelsListMine(service *youtube.Service, part []string) *youtube.ChannelListResponse {
	call := service.Channels.List(part)
	call = call.Mine(true)
	response, err := call.Do()
	handleError(err, "")
	return response
}

func HasLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func GetData(c *gin.Context) {
	config, token, err := getConfigAndToken()
	if err != nil {
		return
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		return
	}

	call := service.Captions.List([]string{"snippet"}, "QSRPKGtjq5A")
	result, err := call.Do()

	// var isLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

	var sb strings.Builder

	for _, val := range result.Items {
		call := service.Captions.Download(val.Id)
		result, err := call.Download()

		if err != nil {
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(result.Body)
		newStr := buf.String()

		for _, line := range strings.Split(newStr, "\n") {
			if HasLetter(line) {
				sb.WriteString(line)
				sb.WriteRune(' ')
			}
		}

		fmt.Println(sb.String())
	}

	response := channelsListMine(service, []string{"contentDetails"})

	for _, channel := range response.Items {
		playlistId := channel.ContentDetails.RelatedPlaylists.Uploads

		// Print the playlist ID for the list of uploaded videos.
		fmt.Printf("Videos in list %s\r\n", playlistId)

		nextPageToken := ""
		for {
			// Retrieve next set of items in the playlist.
			playlistResponse := playlistItemsList(service, []string{"snippet"}, playlistId, nextPageToken)

			for _, playlistItem := range playlistResponse.Items {
				title := playlistItem.Snippet.Title
				videoId := playlistItem.Snippet.ResourceId.VideoId
				fmt.Printf("%v, (%v)\r\n", title, videoId)
			}

			// Set the token to retrieve the next page of results
			// or exit the loop if all results have been retrieved.
			nextPageToken = playlistResponse.NextPageToken
			if nextPageToken == "" {
				break
			}
			fmt.Println()
		}
	}

	c.JSON(200, gin.H{})
	return
}

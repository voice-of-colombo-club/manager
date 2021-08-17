package recording

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/youtube/v3"
)

var tokenMu sync.Mutex

func getConfigAndToken() (*oauth2.Config, *oauth2.Token, error) {
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Unable to read client secret json")
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope, youtube.YoutubeForceSslScope, youtube.YoutubeReadonlyScope, youtube.YoutubepartnerScope)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to create config from json")
	}

	token, err := getTokenFromConfig(config)
	return config, token, err
}

// Retrieve a token, saves the token, then returns the generated client.
func getTokenFromConfig(config *oauth2.Config) (*oauth2.Token, error) {

	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tokenMu.Lock()
		defer tokenMu.Unlock()

		tok, err = getTokenFromWeb(config)
		if err != nil {
			return tok, err
		}
		err = saveToken(tokFile, tok)
	}
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	var authCode = "4/0AX4XfWgNrtEfwk8u9NWRO-KGBlq2ZGazHb4T3TEtR7RkiHSpvjxHGl-Vn7zixT7Hk55xdA"
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		// log.Fatalf("Unable to retrieve token from web: %v", err)
		return tok, err
	}
	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		// log.Fatalf("Unable to cache oauth token: %v", err)
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)

	return nil
}

func printAuthUrl() {
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		fmt.Println(err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope, drive.DriveAppdataScope, drive.DriveFileScope, drive.DriveMetadataReadonlyScope, drive.DriveMetadataScope, drive.DrivePhotosReadonlyScope, drive.DriveReadonlyScope, drive.DriveScope, drive.DriveScriptsScope)
	if err != nil {
		fmt.Println(err)
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println(authURL)
}

type Point struct {
}

func (*Point) doSomething() {

}

func (p *Point) ScaleBy(factor float64) {
}

func abcd(e string) {
	p := Point{}
	(&p).ScaleBy(2)
	fmt.Println(p) // "{2, 4}"
}

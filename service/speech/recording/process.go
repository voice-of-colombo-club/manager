package recording

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/voice-of-colombo/service/database"
	"github.com/voice-of-colombo/service/graph/model"
	"github.com/voice-of-colombo/service/speech"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var mu sync.Mutex

type zoomRecording struct {
	Id             int    `json:"id"`
	Topic          string `json:"topic"`
	TotalSize      int    `json:"total_size"`
	Duration       int    `json:"duration"`
	RecordingFiles []struct {
		FileType      string `json:"file_type"`
		FileExtension string `json:"file_extension"`
		Status        string `json:"status"`
		ShareUrl      string `json:"share_url"`
		DownloadUrl   string `json:"download_url"`
		RecordingType string `json:"recording_type"`
	} `json:"recording_files"`
}

type recordingDownloadResult struct {
	fileExtension string
	downloadUrl   string
	downloadBytes *[]byte
	recordingType string
}

func ProcessRecording(recordingData *model.RecordingProcessInput) error {
	meetingId := recordingData.ZoomMeetingID //"86574878695"
	errGroup, _ := errgroup.WithContext(context.Background())

	tokenString, err := getZoomJwt()
	if err != nil {
		return errors.Wrapf(err, "Error getting zoom jwt")
	}

	downloadData, err := getRecordingsForMeeting(meetingId, tokenString)
	if err != nil {
		return errors.Wrapf(err, "Error getting recordings list for meeting %s", meetingId)
	}
	log.Printf("successfully downloaded data for %s", meetingId)

	primaryRecordingBytes := downloadData["RECORDING"].downloadBytes

	errGroup.Go(func() error {
		err = uploadToDrive(downloadData, meetingId)
		if err != nil {
			return errors.Wrapf(err, "upload recordings to drive failed for %s", meetingId)
		}
		log.Printf("successfully uploaded data to drive for %s", meetingId)
		return nil
	})

	for _, speechEntry := range recordingData.Speeches {
		speechEntry := speechEntry
		errGroup.Go(func() error {
			return processIndividualSpeechRecording(speechEntry, primaryRecordingBytes)
		})
	}

	err = errGroup.Wait()
	if err != nil {
		return errors.Wrapf(err, "errors occurred while processing individual speeches for upload")
	}

	log.Printf("Successfully processed meeting recordings")
	return nil
}

func processIndividualSpeechRecording(speechEntry *model.RecordingSpeechInput, recording *[]byte) error {
	log.Printf("Processing speech for uploaded data to drive for %s", *speechEntry.SpeechID)
	recordingBytes, err := processRecording(recording, speechEntry)

	if err != nil {
		return errors.Wrapf(err, "Unable to process recording for %s", *speechEntry.SpeechID)
	}

	saveSpeech, err := speech.GetSpeechById(*speechEntry.SpeechID)
	if err != nil {
		return errors.Wrapf(err, "Unable to retrieve saved speech for %s", *speechEntry.SpeechID)
	}

	vidId, err := uploadToYoutube(bytes.NewReader(recordingBytes), saveSpeech.Title)
	if err != nil {
		return errors.Wrapf(err, "Unable to upload youtube video for %s", *speechEntry.SpeechID)
	}

	saveSpeech.SpeechLink = vidId
	result := database.Db.Save(saveSpeech)

	if result.Error != nil {
		return errors.Wrapf(err, "Error when db saving video link %s for speech %s", vidId, *speechEntry.SpeechID)
	}

	log.Printf("Successfully uploaded video to %s", *vidId)
	return nil
}

func getCommand(speech *model.RecordingSpeechInput, filePath string) *exec.Cmd {
	commandParams := []string{
		"-loglevel", "repeat+level+error",
		"-analyzeduration", "100M", "-probesize", "100M",
		"-i", "pipe:0",
		"-ss", speech.From,
		"-to", speech.To,
		"-movflags", "frag_keyframe+empty_moov",
		"-f", "mp4",
	}

	if speech.GalleryCrop {
		commandParams = append(commandParams, "-filter:v", "crop=320:180:iw-320:0")
	}

	commandParams = append(commandParams, filePath)

	cmd := exec.Command("ffmpeg", commandParams...)
	cmd.Stderr = os.Stderr // bind log stream to stderr

	return cmd
}

func processRecording(recording *[]byte, speech *model.RecordingSpeechInput) ([]byte, error) {
	log.Printf("Starting youtube upload process")

	// We cannot share the stdout and stdin (as far as I am aware, and had a couple issues when trying to)
	// So we compromise by locking the stdin and without using stdout saving to a file, and releasing the lock after stdin has been read
	tempFileName := uuid.New().String()
	tempFilePath := fmt.Sprintf("/temp/videos/%s.mp4", tempFileName)

	// Wrapped in a function as the function context is mutex locked
	// Stdin cannot be shared
	cmd, err := func() (*exec.Cmd, error) {
		mu.Lock()
		defer mu.Unlock()

		cmd := getCommand(speech, tempFilePath)
		stdin, err := cmd.StdinPipe() // Open stdin pipe
		if err != nil {
			return nil, errors.Wrapf(err, "error opening stdin pipe for %s", *speech.SpeechID)
		}
		defer stdin.Close()

		err = cmd.Start() // Start a process on another goroutine
		if err != nil {
			return nil, errors.Wrapf(err, "error starting the ffmpeg command for %s", *speech.SpeechID)
		}

		_, err = stdin.Write(*recording)
		if !errors.Is(err, syscall.EPIPE) {
			return nil, errors.Wrapf(err, "error writing to the stdin pipe for %s", *speech.SpeechID)
		}

		return cmd, nil
	}()

	if err != nil {
		return nil, err
	}

	// Wait till fmpeg is done processing
	err = cmd.Wait()
	if err != nil {
		return nil, errors.Wrapf(err, "error waiting for ffmpeg to complete processing for %s", *speech.SpeechID)
	}

	tempFileBytes, err := os.ReadFile(tempFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading temp processed file for %s", *speech.SpeechID)
	}

	err = os.Remove(tempFilePath)
	if err != nil {
		log.Printf("Ignoring, unable to clear temp process file for %s", *speech.SpeechID)
	}

	return tempFileBytes, nil
}

// NOTE : This function mutates the map to save memory
func downloadFilesForRecording(downloadurls *map[string]*recordingDownloadResult) (map[string]*[]byte, error) {
	errs, _ := errgroup.WithContext(context.Background())

	downloads := make(map[string]*[]byte)

	// TODO, as improvement we can work on terminating all go routines from the parent upon one error
	// however as we run this only once per zoom recording by an admin, we did not go to that length
	for key, downloadResult := range *downloadurls {
		key := key
		downloadResult := downloadResult

		errs.Go(func() error {
			log.Printf("Initiating download for %s", downloadResult.downloadUrl)

			req, err := http.NewRequest("GET", downloadResult.downloadUrl, nil)
			if err != nil {
				return errors.Wrapf(err, "Eror initializing new request for %s", downloadResult.downloadUrl)
			}

			resp, err := (&http.Client{}).Do(req)
			if err != nil {
				return errors.Wrapf(err, "Error downloading recording for %s", downloadResult.downloadUrl)
			}

			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrapf(err, "An error occurred when converting to bytes %s", key)
			}

			downloads[key] = &respBytes

			log.Printf("Downloaded completed for %s", downloadResult.downloadUrl)
			return nil
		})
	}

	err := errs.Wait()
	return downloads, err
}

func getZoomJwt() (string, error) {
	zoomApiKey := os.Getenv("ZOOM_API_KEY")
	zoomSecret := os.Getenv("ZOOM_SECRET")

	claimDetails := map[string]interface{}{
		"iss": zoomApiKey,
		"exp": "1628698543000", //TODO : Check, dynamic expiry value?
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claimDetails))
	return token.SignedString([]byte(zoomSecret))
}

func uploadToYoutube(r io.Reader, title string) (*string, error) {
	config, token, err := getConfigAndToken()
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to get token for upload")
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating youtube service")
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: title,
			CategoryId:  "22",
			Tags:        strings.Split("test,ee", ","),
		},
		Status: &youtube.VideoStatus{PrivacyStatus: "unlisted"},
	}

	call := service.Videos.Insert([]string{"snippet,status"}, upload)
	if err != nil {
		return nil, errors.Wrapf(err, "error calling video insert") // TODO : find context
	}

	response, err := call.Media(r, googleapi.ContentType("video/mp4")).Do()

	if err != nil {
		return nil, errors.Wrapf(err, "Error when uploading video")
	}

	log.Printf("Upload successful! Video ID: %v", response.Id)
	return &response.Id, nil
}

func getRecordingsForMeeting(meetingId string, token string) (map[string]*recordingDownloadResult, error) {
	zoomMeetingRecordingUrl := fmt.Sprintf("https://api.zoom.us/v2/meetings/%v/recordings", meetingId)
	req, err := http.NewRequest("GET", zoomMeetingRecordingUrl, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating the request for %s", zoomMeetingRecordingUrl)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing request for %s", zoomMeetingRecordingUrl)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading response body for %s", zoomMeetingRecordingUrl)
	}

	results := zoomRecording{}
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, errors.Wrapf(err, "Error converting to json for %s", zoomMeetingRecordingUrl)
	}

	downloadMap := make(map[string]*recordingDownloadResult)

	for _, recEntry := range results.RecordingFiles {
		if recEntry.Status != "completed" {
			return nil, errors.Errorf("Rec status was not completed for %s", recEntry.DownloadUrl)
		}

		key, err := func() (string, error) {
			if recEntry.FileExtension == "MP4" {
				if _, ok := downloadMap["RECORDING"]; ok {
					return "", errors.Errorf("There have been multiple video files in the recording for %s", recEntry.DownloadUrl)
				}
				return "RECORDING", nil
			}
			return recEntry.RecordingType, nil
		}()

		if err != nil {
			return nil, err
		}

		downloadUrlWithToken := fmt.Sprintf("%s?access_token=%s", recEntry.DownloadUrl, token)
		downloadMap[key] = &recordingDownloadResult{recEntry.FileExtension, downloadUrlWithToken, nil, recEntry.RecordingType}
	}

	downloads, err := downloadFilesForRecording(&downloadMap)
	for key, value := range downloadMap {
		value.downloadBytes = downloads[key]
	}
	return downloadMap, err
}

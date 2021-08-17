package recording

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func uploadToDrive(downloads map[string]*recordingDownloadResult, meetingId string) error {
	config, token, err := getConfigAndToken()
	if err != nil {
		return errors.Wrapf(err, "unable to get config or token")
	}

	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		return errors.Wrapf(err, "Unable to retrieve Drive client")
	}

	// TODO : Can we share context? what happens then?
	errGroup, _ := errgroup.WithContext(context.Background())
	for _, downloadResult := range downloads {
		downloadResult := downloadResult

		errGroup.Go(func() error {
			fileName := fmt.Sprintf("%s-%s.%s", meetingId, downloadResult.recordingType, downloadResult.fileExtension)
			upf := srv.Files.Create(&drive.File{Name: fileName}).
				Media(bytes.NewReader(*downloadResult.downloadBytes))

			_, err = upf.Do()
			if err != nil {
				return errors.Wrapf(err, "Unable to upload recording to drive for %s", fileName)
			}

			log.Printf("Successfully uploaded %s", fileName)
			return nil
		})
	}

	return errGroup.Wait()
}

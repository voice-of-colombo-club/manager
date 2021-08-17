package speech

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/voice-of-colombo/service/database"
	"github.com/voice-of-colombo/service/graph/model"
)

func GetSpeechById(id string) (*model.Speech, error) {
	var speech *Speech
	result := database.Db.First(&speech, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}
	return speech.ToGraphQL(), nil
}

func SearchSpeeches(searchCriteria *model.SpeechSearchCriteria) ([]*model.Speech, error) {
	var speeches []*Speech
	queryBuider := database.Db.Limit(*searchCriteria.Limit).Offset(*searchCriteria.Offset)

	if searchCriteria.ID != nil {
		queryBuider = queryBuider.Where("id = ?", *searchCriteria.ID)
	}

	if searchCriteria.UserID != nil {
		queryBuider = queryBuider.Where("user_id = ?", *searchCriteria.UserID)
	}

	// TODO : A bit confusing, to figure out a better alternative
	if len(searchCriteria.SpeechTypes) > 0 {
		queryBuider = queryBuider.Where("speech_type IN ?", searchCriteria.SpeechTypes)
	}

	result := queryBuider.Find(&speeches)

	if result.Error != nil {
		return nil, result.Error
	}
	return SpeechListToGqraphQL(speeches), nil
}

func SaveSpeech(input model.SaveSpeech) (*model.Speech, error) {
	saveSpeech := (&Speech{})
	saveSpeech.LoadSaveSpeechInput(input)

	result := database.Db.Save(saveSpeech)

	if result.Error != nil {
		return nil, result.Error
	}
	return saveSpeech.ToGraphQL(), nil
}

func GetSpeechesByParentId(parentSpeechId string) ([]*model.Speech, error) {
	var speeches []*Speech
	result := database.Db.Where("parent_speech_id = ?", parentSpeechId).Find(&speeches)

	if result.Error != nil {
		return nil, result.Error
	}

	return SpeechListToGqraphQL(speeches), result.Error
}

func DownloadEvaluation(c *gin.Context) {
	speechId := c.Param("speechId")
	speech, err := GetSpeechById(speechId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrapf(err, "Unable to load speech %s", speechId))
		return
	}

	_, fileName := filepath.Split(*speech.EvaluationFilePath)
	c.FileAttachment(*speech.EvaluationFilePath, fileName)
}

func SaveEvaluation(c *gin.Context) {
	speechId := c.Param("speechId")

	speech, err := GetSpeechById(speechId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Wrapf(err, "Unable to load speech %s", speechId))
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	extension := filepath.Ext(header.Filename)
	fmt.Println(extension)
	if !strings.EqualFold(extension, ".pdf") {
		c.String(http.StatusBadRequest, "File type should be pdf")
		return
	}

	if speech.SpeechType != "Evaluation" {
		c.String(http.StatusBadRequest, "Speech type needs to be evaluation")
		return
	}

	path := fmt.Sprintf("%s-%s.pdf", speechId, uuid.NewV4().String())
	out, err := os.Create(path)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	speech.EvaluationFilePath = &path
	result := database.Db.Save(speech)
	if result.Error != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"filepath": path})
}

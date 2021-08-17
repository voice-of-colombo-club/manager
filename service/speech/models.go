package speech

import (
	uuid "github.com/satori/go.uuid"
	"github.com/voice-of-colombo/service/graph/model"
	"gorm.io/gorm"
)

type Speech struct {
	gorm.Model
	ID                 string  `json:"id"`
	MeetingID          *string `json:"meetingId"`
	SpeechType         string  `json:"speechType"`
	Timestamp          *int    `json:"timestamp"`
	Location           *string `json:"location"`
	Project            *string `json:"project"`
	UserID             string  `json:"userId"`
	SpeechLink         *string `json:"speechLink"`
	Title              string  `json:"title"`
	ParentSpeechID     *string `json:"parentSpeechId"`
	ActualTime         string  `json:"actualTime"`
	ExpectedTime       string  `json:"expectedTime"`
	EvaluationFilePath *string `json:"evaluationFile"`
}

func (s *Speech) ToGraphQL() *model.Speech {
	gormSpeech := model.Speech{
		ID:                 s.ID,
		MeetingID:          s.MeetingID,
		SpeechType:         s.SpeechType,
		Timestamp:          s.Timestamp,
		Location:           s.Location,
		Project:            s.Project,
		UserID:             s.UserID,
		SpeechLink:         s.SpeechLink,
		Title:              s.Title,
		ParentSpeechID:     s.ParentSpeechID,
		ActualTime:         s.ActualTime,
		ExpectedTime:       s.ExpectedTime,
		EvaluationFilePath: s.EvaluationFilePath,
	}
	return &gormSpeech
}

func (s *Speech) LoadSaveSpeechInput(input model.SaveSpeech) {
	id := func() string {
		if input.ID == nil {
			return uuid.NewV4().String()
		}
		return *input.ID
	}()

	s.SpeechType = input.SpeechType
	s.MeetingID = input.MeetingID
	s.UserID = input.UserID
	s.Timestamp = input.Timestamp
	s.Location = input.Location
	s.Project = input.Project
	s = &Speech{
		ID:         id,
		SpeechType: input.SpeechType,
		MeetingID:  input.MeetingID,
		UserID:     input.UserID,
		Timestamp:  input.Timestamp,
		Location:   input.Location,
		Project:    input.Project,
		SpeechLink: input.Project,
	}
}

func SpeechToModel(s *model.Speech) *Speech {
	gqlSpeech := Speech{
		ID:                 s.ID,
		MeetingID:          s.MeetingID,
		SpeechType:         s.SpeechType,
		Timestamp:          s.Timestamp,
		Location:           s.Location,
		Project:            s.Project,
		UserID:             s.UserID,
		SpeechLink:         s.SpeechLink,
		Title:              s.Title,
		ParentSpeechID:     s.ParentSpeechID,
		ActualTime:         s.ActualTime,
		ExpectedTime:       s.ExpectedTime,
		EvaluationFilePath: s.EvaluationFilePath,
	}
	return &gqlSpeech
}

func SpeechListToGqraphQL(speeches []*Speech) []*model.Speech {
	gqlSpeeches := make([]*model.Speech, 0)
	for _, speech := range speeches {
		gqlSpeeches = append(gqlSpeeches, speech.ToGraphQL())
	}
	return gqlSpeeches
}

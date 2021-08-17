package meeting

import (
	uuid "github.com/satori/go.uuid"
	"github.com/voice-of-colombo/service/graph/model"
	"gorm.io/gorm"
)

type Meeting struct {
	gorm.Model
	ID                   string  `json:"id"`
	Number               int     `json:"number"`
	Theme                string  `json:"theme"`
	Location             string  `json:"location"`
	JointMeetingClubName *string `json:"jointMeetingClubName"`
	IsAreaMeeting        bool    `json:"isAreaMeeting"`
	Timestamp            int     `json:"timestamp"`
}

func (m *Meeting) ToGraphQL() *model.Meeting {
	gormSpeech := model.Meeting{
		ID:                   m.ID,
		Number:               m.Number,
		Theme:                m.Theme,
		Location:             m.Location,
		JointMeetingClubName: m.JointMeetingClubName,
		IsAreaMeeting:        m.IsAreaMeeting,
		Timestamp:            m.Timestamp,
	}
	return &gormSpeech
}

func (m *Meeting) LoadSaveMeetingInput(input model.CreateMeeting) {
	id := func() string {
		if input.ID == nil {
			return uuid.NewV4().String()
		}
		return *input.ID
	}()

	jointMeetingClubName := func() string {
		if input.JointMeetingClubName == nil {
			return ""
		}
		return *input.JointMeetingClubName
	}()

	m.ID = id
	m.Number = input.Number
	m.Theme = input.Theme
	m.Location = input.Location
	m.JointMeetingClubName = &jointMeetingClubName
	m.IsAreaMeeting = input.IsAreaMeeting
	m.Timestamp = input.Timestamp
}

func MeetingToModel(m *model.Meeting) *Meeting {
	gqlMeeting := Meeting{
		ID:                   m.ID,
		Number:               m.Number,
		Theme:                m.Theme,
		Location:             m.Location,
		JointMeetingClubName: m.JointMeetingClubName,
		IsAreaMeeting:        m.IsAreaMeeting,
		Timestamp:            m.Timestamp,
	}
	return &gqlMeeting
}

func MeetingListToGqraphQL(meetings []*Meeting) []*model.Meeting {
	gqlMeetings := make([]*model.Meeting, 0)
	for _, meeting := range meetings {
		gqlMeetings = append(gqlMeetings, meeting.ToGraphQL())
	}
	return gqlMeetings
}

package meeting

import (
	"github.com/voice-of-colombo/service/database"
	"github.com/voice-of-colombo/service/graph/model"
)

func GetMeetingById(id string) (*model.Meeting, error) {
	var meeting *Meeting
	result := database.Db.First(&meeting, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}
	return meeting.ToGraphQL(), nil
}

func SearchMeetings(searchCriteria *model.MeetingSearchCriteria) ([]*model.Meeting, error) {
	var meetings []*Meeting
	queryBuider := database.Db.Limit(*searchCriteria.Limit).Offset(*searchCriteria.Offset)
	result := queryBuider.Find(&meetings)

	if result.Error != nil {
		return nil, result.Error
	}
	return MeetingListToGqraphQL(meetings), nil
}

func CreateMeeting(input model.CreateMeeting) (*model.Meeting, error) {
	newMeeting := &Meeting{}
	newMeeting.LoadSaveMeetingInput(input)

	result := database.Db.Save(newMeeting)

	if result.Error != nil {
		return nil, result.Error
	}
	return newMeeting.ToGraphQL(), nil
}

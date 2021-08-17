package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/voice-of-colombo/service/graph/generated"
	"github.com/voice-of-colombo/service/graph/model"
	"github.com/voice-of-colombo/service/meeting"
	"github.com/voice-of-colombo/service/speech"
	"github.com/voice-of-colombo/service/user"
)

func (r *speechResolver) User(ctx context.Context, obj *model.Speech) (*model.User, error) {
	limit := 2
	searchCriteria := model.UserSearchCriteria{
		ID:     &obj.UserID,
		Offset: new(int),
		Limit:  &limit,
	}
	users, err := user.SearchUsers(searchCriteria)
	return users[0], err
}

func (r *speechResolver) Meeting(ctx context.Context, obj *model.Speech) (*model.Meeting, error) {
	if obj.MeetingID == nil {
		return &model.Meeting{}, nil
	}

	return meeting.GetMeetingById(*obj.MeetingID)
}

func (r *speechResolver) ParentSpeech(ctx context.Context, obj *model.Speech) (*model.Speech, error) {
	if obj.ParentSpeechID == nil {
		return &model.Speech{}, nil
	}

	return speech.GetSpeechById(*obj.ParentSpeechID)
}

func (r *speechResolver) AhCounts(ctx context.Context, obj *model.Speech) (*model.AhCounts, error) {
	return speech.GetAhCountsBySpeechId(obj.ID)
}

func (r *speechResolver) ChildSpeeches(ctx context.Context, obj *model.Speech) ([]*model.Speech, error) {
	return speech.GetSpeechesByParentId(obj.ID)
}

// Speech returns generated.SpeechResolver implementation.
func (r *Resolver) Speech() generated.SpeechResolver { return &speechResolver{r} }

type speechResolver struct{ *Resolver }

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/voice-of-colombo/service/graph/generated"
	"github.com/voice-of-colombo/service/graph/model"
	"github.com/voice-of-colombo/service/meeting"
	"github.com/voice-of-colombo/service/speech"
	"github.com/voice-of-colombo/service/speech/recording"
	"github.com/voice-of-colombo/service/user"
)

func (r *mutationResolver) CreateMeeting(ctx context.Context, input model.CreateMeeting) (*model.Meeting, error) {
	return meeting.CreateMeeting(input)
}

func (r *mutationResolver) SaveSpeech(ctx context.Context, input model.SaveSpeech) (*model.Speech, error) {
	return speech.SaveSpeech(input)
}

func (r *mutationResolver) ProcessRecording(ctx context.Context, input model.RecordingProcessInput) (bool, error) {
	err := recording.ProcessRecording(&input)

	if err != nil {
		log.Printf("Errors occurred when processing the recording %v", err)
		return false, err
	}
	return true, nil
}

func (r *queryResolver) Meetings(ctx context.Context, searchCriteria *model.MeetingSearchCriteria) ([]*model.Meeting, error) {
	return meeting.SearchMeetings(searchCriteria)
}

func (r *queryResolver) Meeting(ctx context.Context, id string) (*model.Meeting, error) {
	return meeting.GetMeetingById(id)
}

func (r *queryResolver) Users(ctx context.Context, searchCriteria model.UserSearchCriteria) ([]*model.User, error) {
	return user.SearchUsers(searchCriteria)
}

func (r *queryResolver) Speeches(ctx context.Context, searchCriteria *model.SpeechSearchCriteria) ([]*model.Speech, error) {
	return speech.SearchSpeeches(searchCriteria)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

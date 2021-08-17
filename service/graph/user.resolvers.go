package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/voice-of-colombo/service/graph/generated"
	"github.com/voice-of-colombo/service/graph/model"
	"github.com/voice-of-colombo/service/user"
)

func (r *userResolver) Clubs(ctx context.Context, obj *model.User) ([]*model.Club, error) {
	return user.GetClubsForUser(obj.ID)
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

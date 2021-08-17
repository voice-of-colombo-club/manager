package user

import (
	"github.com/voice-of-colombo/service/graph/model"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string
	IsAdmin   bool `json:"isAdmin" gorm:"default:false"`
}

func (u *User) ToGraphQL() *model.User {
	gormUser := model.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
	}
	return &gormUser
}

func UserListToGqraphQL(users []*User) []*model.User {
	gqlUsers := make([]*model.User, 0)
	for _, user := range users {
		gqlUsers = append(gqlUsers, user.ToGraphQL())
	}
	return gqlUsers
}

type Club struct {
	gorm.Model
	ID       string `json:"id"`
	ClubName string `json:"clubName"`
}

func (c *Club) ToGraphQL() *model.Club {
	gormClub := model.Club{
		ID:       c.ID,
		ClubName: c.ClubName,
	}
	return &gormClub
}

type ClubUserMapping struct {
	gorm.Model
	ClubId string
	UserId string
}

func ClubListToGqraphQL(clubs []*Club) []*model.Club {
	gqlClubs := make([]*model.Club, 0)
	for _, club := range clubs {
		gqlClubs = append(gqlClubs, club.ToGraphQL())
	}
	return gqlClubs
}

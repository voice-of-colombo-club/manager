package model

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"isAdmin" gorm:"default:false"`
}

type Club struct {
	ID       string `json:"id"`
	ClubName string `json:"clubName"`
}

type ClubUserMapping struct {
	ClubId string
	UserId string
}

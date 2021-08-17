package model

type Meeting struct {
	ID                   string  `json:"id"`
	Number               int     `json:"number"`
	Theme                string  `json:"theme"`
	Location             string  `json:"location"`
	JointMeetingClubName *string `json:"jointMeetingClubName"`
	IsAreaMeeting        bool    `json:"isAreaMeeting"`
	Timestamp            int     `json:"timestamp"`
}

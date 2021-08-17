package model

type Speech struct {
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

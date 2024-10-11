package models

type SmsToUrlForwarderMessageDTO struct {
	From         string `json:"from"`
	Text         string `json:"text"`
	SentAtTs     int64  `json:"sentStamp"`
	ReceivedAtTs int64  `json:"receivedStamp"`
	Sim          string `json:"sim"`
}

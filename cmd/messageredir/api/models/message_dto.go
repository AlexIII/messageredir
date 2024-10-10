package models

type SmsToUrlForwarderMessageDTO struct {
	From     string `json:"from"`
	Text     string `json:"text"`
	Sent     string `json:"sentStamp"`
	Received string `json:"receivedStamp"`
	Sim      string `json:"sim"`
}

package models

type TelegramMessageIn struct {
	ChatId   int64  // Sender's chat ID
	Username string // Sender's username
	Text     string
	Command  string
}

type TelegramMessageOut struct {
	ChatId int64 // Recipient's chat ID
	Text   string
}

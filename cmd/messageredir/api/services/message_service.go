package services

type MessageService interface {
	SendStr(chatId int64, message string) error
}

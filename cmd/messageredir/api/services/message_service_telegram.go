package services

type OutMessage struct {
	ChatId  int64
	Message string
}

type MessageServiceTelegram struct {
	outMessages chan OutMessage
}

func NewMessageServiceTelegram() MessageServiceTelegram {
	return MessageServiceTelegram{make(chan OutMessage, 10)}
}

func (service MessageServiceTelegram) SendStr(chatId int64, message string) error {
	service.outMessages <- OutMessage{chatId, message}
	return nil
}

func (service MessageServiceTelegram) GetOutMessages() <-chan OutMessage {
	return service.outMessages
}

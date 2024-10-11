package services

import "messageredir/cmd/messageredir/services/models"

type TelegramService interface {
	Send(message models.TelegramMessageOut) error
	GetReceiveChan() <-chan models.TelegramMessageIn
}

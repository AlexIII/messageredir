package services

import (
	"log"
	"messageredir/cmd/messageredir/services/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegram struct {
	outMsg chan models.TelegramMessageOut
	inMsg  chan models.TelegramMessageIn
}

func StartTelegramService(botToken string) TelegramService {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	// bot.Debug = true
	log.Printf("Authorized on telegram account %s", bot.Self.UserName)

	tg := telegram{
		outMsg: make(chan models.TelegramMessageOut, 10),
		inMsg:  make(chan models.TelegramMessageIn),
	}

	go func() {
		updConf := tgbotapi.NewUpdate(0)
		updConf.Timeout = 60
		updates := bot.GetUpdatesChan(updConf)
		for {
			select {
			// Send outgoing message
			case msg := <-tg.outMsg:
				tgMsg := tgbotapi.NewMessage(msg.ChatId, msg.Text)
				tgMsg.DisableWebPagePreview = true
				bot.Send(tgMsg)
			// Queue incoming message
			case update := <-updates:
				if update.Message != nil {
					tg.inMsg <- models.TelegramMessageIn{
						ChatId:   update.Message.Chat.ID,
						Username: update.Message.Chat.UserName,
						Text:     update.Message.Text,
						Command:  update.Message.Command(),
					}
				}
			}
		}
	}()

	return tg
}

func (tg telegram) Send(message models.TelegramMessageOut) error {
	tg.outMsg <- message
	return nil
}

func (tg telegram) GetReceiveChan() <-chan models.TelegramMessageIn {
	return tg.inMsg
}

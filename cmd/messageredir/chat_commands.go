package main

import (
	"log"
	"messageredir/cmd/messageredir/services/models"
)

func (cmd App) start(message models.TelegramMessageIn) {
	user := cmd.dbRepo.GetOrCreateUser(message.ChatId, message.Username, cmd.config.UserTokenLength)
	cmd.tgService.Send(models.TelegramMessageOut{
		ChatId: user.ChatId,
		Text:   "You are good to go!\nYour token: " + user.Token,
	})
}

func (cmd App) end(message models.TelegramMessageIn) {
	cmd.dbRepo.DeleteUser(message.ChatId)
	cmd.tgService.Send(models.TelegramMessageOut{
		ChatId: message.ChatId,
		Text:   "You were erased from the system. Goodbye!",
	})
}

func (cmd App) Process(message models.TelegramMessageIn) bool {
	cmdMap := map[string]func(message models.TelegramMessageIn){
		"start": cmd.start,
		"end":   cmd.end,
	}
	if run, found := cmdMap[message.Command]; found {
		log.Println("Processing command", message.Command, message.ChatId)
		run(message)
		return true
	}
	return false
}

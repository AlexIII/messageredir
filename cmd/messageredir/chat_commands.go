package main

import (
	"log"
	"messageredir/cmd/messageredir/services/models"
)

func (ctx *App) start(message models.TelegramMessageIn) {
	user := ctx.db.GetOrCreateUser(message.ChatId, message.Username, ctx.config.UserTokenLength)
	ctx.telegram.Send(models.TelegramMessageOut{
		ChatId: user.ChatId,
		Text:   "You are good to go!\nYour token: " + user.Token,
	})
}

func (ctx *App) end(message models.TelegramMessageIn) {
	ctx.db.DeleteUser(message.ChatId)
	ctx.telegram.Send(models.TelegramMessageOut{
		ChatId: message.ChatId,
		Text:   "You were erased from the system. Goodbye!",
	})
}

var commandMap = map[string]func(cmd *App, message models.TelegramMessageIn){
	"start": (*App).start,
	"end":   (*App).end,
}

func (ctx *App) Process(message models.TelegramMessageIn) bool {
	if run, found := commandMap[message.Command]; found {
		log.Println("Processing command", message.Command, "for", message.ChatId)
		run(ctx, message)
		return true
	}
	return false
}

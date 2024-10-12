package main

import (
	"fmt"
	"log"
	"messageredir/cmd/messageredir/services/models"
	"messageredir/cmd/messageredir/strings"
)

func (ctx *App) start(message models.TelegramMessageIn) {
	user := ctx.db.GetOrCreateUser(message.ChatId, message.Username, ctx.config.UserTokenLength)

	ctx.telegram.Send(models.TelegramMessageOut{
		ChatId: user.ChatId,
		Text:   fmt.Sprintf(strings.UserAdded, user.Token),
	})

	// Separate URL message
	protocol := "http"
	if ctx.config.IsTlsEnabled() {
		protocol = "https"
	}
	host := ctx.config.MyHost
	if host == "" {
		host = "<DOMAIN_OR_IP>:<PORT>"
	}
	ctx.telegram.Send(models.TelegramMessageOut{
		ChatId: user.ChatId,
		Text:   fmt.Sprintf("%s://%s/%s/smstourlforwarder", protocol, host, user.Token),
	})
}

func (ctx *App) end(message models.TelegramMessageIn) {
	ctx.db.DeleteUser(message.ChatId)
	ctx.telegram.Send(models.TelegramMessageOut{
		ChatId: message.ChatId,
		Text:   strings.UserRemoved,
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

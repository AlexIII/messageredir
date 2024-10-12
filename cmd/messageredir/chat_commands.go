package main

import (
	"fmt"
	"log"
	sm "messageredir/cmd/messageredir/services/models"
	"messageredir/cmd/messageredir/userstrings"
	"strconv"
	"strings"
)

func (ctx *App) start(message sm.TelegramMessageIn) bool {
	user := ctx.db.GetOrCreateUser(message.ChatId, message.Username, ctx.config.UserTokenLength)

	ctx.telegram.Send(sm.TelegramMessageOut{
		ChatId: user.ChatId,
		Text:   fmt.Sprintf(userstrings.UserAdded, user.Token),
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
	ctx.telegram.Send(sm.TelegramMessageOut{
		ChatId: user.ChatId,
		Text:   fmt.Sprintf("%s://%s/%s/smstourlforwarder", protocol, host, user.Token),
	})

	return true
}

func (ctx *App) end(message sm.TelegramMessageIn) bool {
	ctx.db.DeleteUser(message.ChatId)
	ctx.telegram.Send(sm.TelegramMessageOut{
		ChatId: message.ChatId,
		Text:   userstrings.UserRemoved,
	})

	return true
}

func (ctx *App) conf(message sm.TelegramMessageIn) bool {
	user := ctx.db.GetUserByChatId(message.ChatId)
	if user == nil {
		return false
	}

	// Parse user preferences
	parsedAny := false
	args, _ := strings.CutPrefix(message.Text, "/config")
	for _, elem := range strings.Split(args, ",") {
		kv := strings.Split(elem, ":")
		if len(kv) < 2 {
			continue
		}
		k, v := strings.ToLower(strings.TrimSpace(kv[0])), strings.TrimSpace(kv[1])
		matchedAny := true
		switch k {
		case "utc", "gmt":
			if v, err := strconv.ParseFloat(v, 32); err == nil && v >= -12 && v <= 12 {
				user.Preferences.UtcOffset = int(v * 60)
			}
		case "sim1", "sim2", "sim3", "sim4":
			if user.Preferences.SimNames == nil {
				user.Preferences.SimNames = make(map[string]string)
			}
			user.Preferences.SimNames[k] = v
		default:
			matchedAny = false
		}
		if matchedAny {
			parsedAny = true
		}
	}

	// Send response
	resp := userstrings.UserPreferencesHint
	if parsedAny {
		ctx.db.UpdateUserPreferences(user.ID, user.Preferences)
		resp = userstrings.UserPreferencesUpdated
	}
	ctx.telegram.Send(sm.TelegramMessageOut{
		ChatId: message.ChatId,
		Text:   fmt.Sprintf(resp, user.Preferences),
	})

	return true
}

var commandMap = map[string]func(cmd *App, message sm.TelegramMessageIn) bool{
	"start":  (*App).start,
	"end":    (*App).end,
	"config": (*App).conf,
}

func (ctx *App) Process(message sm.TelegramMessageIn) bool {
	if run, found := commandMap[message.Command]; found {
		log.Println("Processing command", message.Command, "for", message.ChatId)
		return run(ctx, message)
	}
	return false
}

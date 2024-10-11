package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"messageredir/cmd/messageredir/api/middleware"
	am "messageredir/cmd/messageredir/api/models"
	"messageredir/cmd/messageredir/config"
	db "messageredir/cmd/messageredir/db/models"
	"messageredir/cmd/messageredir/db/repo"
	"messageredir/cmd/messageredir/services"
	sm "messageredir/cmd/messageredir/services/models"
	"messageredir/cmd/messageredir/strings"
	"net/http"
)

type MessageController struct {
	config   *config.Config
	db       repo.DbRepo
	telegram services.TelegramService
}

func NewMessageController(config *config.Config, db repo.DbRepo, telegram services.TelegramService) MessageController {
	return MessageController{config, db, telegram}
}

func (ctx MessageController) SmsToUrlForwarder(w http.ResponseWriter, r *http.Request) {
	var message am.SmsToUrlForwarderMessageDTO
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		log.Println("Error parsing request:", err)
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(middleware.UserKey).(*db.User)
	if !ok {
		panic("User not found in context")
	}

	if ctx.config.LogUserMessages {
		log.Printf("New message: %s %+v", user.Username, message)
	}
	ctx.db.UpdateUserStats(user.ID, repo.UpdateUserStats{MessageRedir: true})

	log.Println("Pushing new message for user", user.Username)
	ctx.telegram.Send(sm.TelegramMessageOut{
		ChatId: user.ChatId,
		Text:   fmt.Sprintf(strings.MsgRedirFmt, message.From, message.Sent, message.Sim, message.Text),
	})
}

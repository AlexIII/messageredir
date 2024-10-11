package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"messageredir/cmd/messageredir/api/middleware"
	api "messageredir/cmd/messageredir/api/models"
	"messageredir/cmd/messageredir/api/services"
	"messageredir/cmd/messageredir/app"
	db "messageredir/cmd/messageredir/db/models"
	"net/http"
)

type MessageController struct {
	Config         app.Config
	MessageService services.MessageService
}

func NewMessageController(config app.Config, messageService services.MessageService) MessageController {
	return MessageController{config, messageService}
}

func messageToString(message api.SmsToUrlForwarderMessageDTO) string {
	return fmt.Sprintf("From: %s\nDate: %s\nSim: %s\n\n%s", message.From, message.Sent, message.Sim, message.Text)
}

func (ctx MessageController) SmsToUrlForwarder(w http.ResponseWriter, r *http.Request) {
	var message api.SmsToUrlForwarderMessageDTO
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		log.Println("Error parsing request:", err)
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(middleware.UserKey).(*db.User)
	if !ok {
		panic("User not found in context")
	}

	if ctx.Config.LogUserMessages {
		log.Printf("New message: %s %+v", user.Username, message)
	}

	log.Println("Pushing new message for user", user.Username)
	ctx.MessageService.SendStr(user.ChatId, messageToString(message))
}

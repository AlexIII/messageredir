package controllers

import (
	"encoding/json"
	"fmt"
	"io"
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
	"time"
	"unicode"
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var message am.SmsToUrlForwarderMessageDTO
	if err := json.Unmarshal(body, &message); err != nil {
		log.Println("Error parsing request:", err)
		log.Println("Request body:", safeConvert(body))
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
		Text:   fmt.Sprintf(strings.MsgRedirFmt, formatTimestamp(message.SentAtTs), message.From, message.Sim, message.Text),
	})
}

func safeConvert(data []byte) string {
	result := make([]rune, 0, len(data))
	for _, b := range data {
		if unicode.IsPrint(rune(b)) {
			result = append(result, rune(b))
		}
	}
	return string(result)
}

func formatTimestamp(ms int64) string {
	return time.Unix(ms/1000, 0).UTC().Format("2006-01-02 15:04:05 UTC")
}

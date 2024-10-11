package repo

import (
	"messageredir/cmd/messageredir/db/models"
)

type UpdateUserStats struct {
	MessageRedir bool
}

type DbRepo interface {
	DeleteUser(chatId int64)
	GetUserByToken(token string) *models.User
	GetUserByChatId(chatId int64) *models.User
	GetOrCreateUser(chatId int64, username string, generateNewTokenLength int) *models.User
	UpdateUserStats(userId uint, stats UpdateUserStats)
}

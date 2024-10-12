package repo

import (
	"messageredir/cmd/messageredir/db/models"
)

type UpdateUserStats struct {
	MessageRedir bool // Update last redirected message time (User.LastMessageRedirAt)
}

type DbRepo interface {
	DeleteUser(chatId int64)
	GetUserByToken(token string) *models.User
	GetUserByChatId(chatId int64) *models.User
	GetOrCreateUser(chatId int64, username string, generateNewTokenLength int) *models.User
	UpdateUserStats(userId uint, stats UpdateUserStats)
	UpdateUserPreferences(userId uint, preferences models.UserPreferences)
}

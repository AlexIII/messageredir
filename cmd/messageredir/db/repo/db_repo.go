package repo

import (
	"messageredir/cmd/messageredir/db/models"
)

type DbRepo interface {
	DeleteUser(chatId int64)
	GetUserByToken(token string) *models.User
	GetUserByChatId(chatId int64) *models.User
	GetOrCreateUser(chatId int64, username string, generateNewTokenLength int) *models.User
}

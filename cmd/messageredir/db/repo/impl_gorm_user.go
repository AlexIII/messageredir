package repo

import (
	"errors"
	"log"
	"messageredir/cmd/messageredir/accessToken"
	"messageredir/cmd/messageredir/db/models"
	"time"

	"gorm.io/gorm"
)

func (db *DbRepoGorm) DeleteUser(chatId int64) {
	var user models.User
	if err := db.driver.Where("chat_id = ?", chatId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("User with chat_id =", chatId, "not found")
			return // Nothing to do
		} else {
			log.Panic(err)
		}
	}

	log.Println("Deleting User:", user.ChatId, user.Username)
	err := db.driver.Delete(&user).Error
	if err != nil {
		log.Panic(err)
	}
}

func (db *DbRepoGorm) UpdateUserStats(userId uint, stats UpdateUserStats) {
	user := db.getUserBy("id = ?", userId)
	if stats.MessageRedir {
		user.LastMessageRedirAt = time.Now().UTC()
	}
	err := db.driver.Save(user).Error
	if err != nil {
		log.Panic(err)
	}
}

func (db *DbRepoGorm) getUserBy(query interface{}, args ...interface{}) *models.User {
	var user models.User
	if err := db.driver.Where(query, args).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("User with query '", query, "'", args, "not found")
			return nil
		} else {
			log.Panic(err)
		}
	}

	log.Println("Found User:", user.ChatId, user.Username)
	return &user
}

func (db *DbRepoGorm) GetUserByToken(token string) *models.User {
	return db.getUserBy("token = ?", token)
}

func (db *DbRepoGorm) GetUserByChatId(chatId int64) *models.User {
	return db.getUserBy("chat_id = ?", chatId)
}

// Guaranteed to return a user
func (db *DbRepoGorm) GetOrCreateUser(chatId int64, username string, generateNewTokenLength int) *models.User {
	user := db.GetUserByChatId(chatId)

	newToken := func() string {
		tk, err := accessToken.Generate(generateNewTokenLength)
		if err != nil {
			log.Panic(err)
		}
		return tk
	}

	if user != nil {
		if generateNewTokenLength > 0 {
			user.Token = newToken()
			err := db.driver.Save(user).Error
			if err != nil {
				log.Panic(err)
			}
		}
		return user
	}

	// Create new user
	user = &models.User{
		ChatId:   chatId,
		Username: username,
		Token:    newToken(),
	}
	err := db.driver.Create(user).Error
	if err != nil {
		log.Panic(err)
	}

	return user
}

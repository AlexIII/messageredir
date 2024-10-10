package repo

import (
	"errors"
	"log"
	"messageredir/cmd/messageredir/accessToken"
	"messageredir/cmd/messageredir/db/models"

	"gorm.io/gorm"
)

func DeleteUser(db *gorm.DB, chatId int64) {
	var user models.User
	if err := db.Where("chat_id = ?", chatId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("User with chat_id =", chatId, "not found")
			return // Nothing to do
		} else {
			log.Panic(err)
		}
	}

	log.Println("Deleting User:", user.ChatId, user.Username)
	err := db.Delete(&user).Error
	if err != nil {
		log.Panic(err)
	}
}

func getUserBy(db *gorm.DB, query interface{}, args ...interface{}) *models.User {
	var user models.User
	if err := db.Where(query, args).First(&user).Error; err != nil {
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

func GetUserByToken(db *gorm.DB, token string) *models.User {
	return getUserBy(db, "token = ?", token)
}

func GetUserByChatId(db *gorm.DB, chatId int64) *models.User {
	return getUserBy(db, "chat_id = ?", chatId)
}

// Guaranteed to return a user
func GetOrCreateUser(db *gorm.DB, chatId int64, username string, generateNewTokenLength int) *models.User {
	user := GetUserByChatId(db, chatId)

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
			err := db.Save(&user).Error
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
	err := db.Create(&user).Error
	if err != nil {
		log.Panic(err)
	}

	return user
}

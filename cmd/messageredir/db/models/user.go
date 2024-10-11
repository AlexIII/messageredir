package models

import "time"

type User struct {
	ID                 uint  `gorm:"primaryKey"`
	ChatId             int64 `gorm:"uniqueIndex"`
	Username           string
	Token              string    `gorm:"uniqueIndex"`
	LastMessageRedirAt time.Time `gorm:"default:'1970-01-01 00:00:00'"`
}

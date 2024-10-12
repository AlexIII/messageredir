package models

import "time"

type UserPreferences struct {
	UtcOffset int               `json:"utcOffset"` // In minutes
	SimNames  map[string]string `json:"simNames"`  // sim1 -> "AT&T", sim2 -> "Verizon"
}

type User struct {
	ID                 uint  `gorm:"primaryKey"`
	ChatId             int64 `gorm:"uniqueIndex"`
	Username           string
	Token              string          `gorm:"uniqueIndex"`
	LastMessageRedirAt time.Time       `gorm:"default:'1970-01-01 00:00:00'"`
	Preferences        UserPreferences `gorm:"serializer:json"`
}

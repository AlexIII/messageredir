package models

type User struct {
	ID       uint `gorm:"primaryKey"`
	ChatId   int64
	Username string
	Token    string
}

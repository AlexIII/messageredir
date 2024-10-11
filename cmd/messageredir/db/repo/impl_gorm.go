package repo

import (
	"log"
	"messageredir/cmd/messageredir/db/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbRepoGorm struct {
	driver *gorm.DB
}

func NewDbRepoGorm(dbFileName string) DbRepo {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect to database")
	}
	db.AutoMigrate(&models.User{})
	return &DbRepoGorm{db}
}

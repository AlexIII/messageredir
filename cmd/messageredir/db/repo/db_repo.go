package repo

import "gorm.io/gorm"

type DbRepo struct {
	driver *gorm.DB
}

func NewDbRepo(db *gorm.DB) DbRepo {
	return DbRepo{db}
}

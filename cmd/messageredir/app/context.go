package app

import (
	"gorm.io/gorm"
)

type Context struct {
	Config Config
	Db     *gorm.DB
}
